package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ca-risken/common/pkg/prowler"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v3"
)

const (
	REPO_URL    = "https://github.com/prowler-cloud/prowler.git"
	TMP_DIR     = "./tmp"
	PLUGIN_DIR  = "prowler/providers/azure/services"
	PLUGIN_FILE = "../../prowler.yaml"
)

var (
	// parameters
	pluginDir  = ""
	pluginFile = ""
	commitHash = ""
	isNew      = false
)

func init() {
	pluginDir = PLUGIN_DIR
	if os.Getenv("PLUGIN_DIR") != "" {
		pluginDir = os.Getenv("PLUGIN_DIR")
	}

	pluginFile = PLUGIN_FILE
	if os.Getenv("PLUGIN_FILE") != "" {
		pluginFile = os.Getenv("PLUGIN_FILE")
	} else {
		isNew = true
	}

	if os.Getenv("COMMIT_HASH") != "" {
		commitHash = os.Getenv("COMMIT_HASH")
	}
}

func main() {
	// Prowlerの最新プラグインを取得
	remotePlugin, err := getRemotePlugin()
	if err != nil {
		log.Fatalf("Failed to get remote plugin: %v", err)
	}
	currentPlugin := &prowler.ProwlerSetting{}
	if !isNew {
		currentPlugin, err = prowler.LoadProwlerSetting(pluginFile)
		if err != nil {
			log.Fatalf("Failed to load current plugin: %v", err)
		}
	}

	// データ更新（差分）
	err = updatePlugin(currentPlugin, remotePlugin)
	if err != nil {
		log.Fatalf("Failed to update plugin: %v", err)
	}
	fmt.Println("Completed processing and cleaned up temporary files.")
}

func getRemotePlugin() (*prowler.ProwlerSetting, error) {
	err := os.MkdirAll(TMP_DIR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(TMP_DIR)

	// tmpディレクトリにクローン
	repoDir := filepath.Join(TMP_DIR, "prowler")
	_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:      REPO_URL,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to clone repo: %v", err)
	}

	// 指定されたcommit hashにチェックアウト
	if commitHash != "" {
		repo, err := git.PlainOpen(repoDir)
		if err != nil {
			return nil, fmt.Errorf("Failed to open repository: %v", err)
		}
		worktree, err := repo.Worktree()
		if err != nil {
			return nil, fmt.Errorf("Failed to get worktree: %v", err)
		}
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commitHash),
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to checkout commit %s: %v", commitHash, err)
		}
	}

	// プラグインディレクトリを処理 (providers/azure/services)
	setting, err := processServices(filepath.Join(repoDir, pluginDir))
	if err != nil {
		return nil, fmt.Errorf("Failed to process JS files: %v", err)
	}
	return setting, nil
}

func processServices(baseDir string) (*prowler.ProwlerSetting, error) {
	setting := prowler.ProwlerSetting{
		SpecificPluginSetting: map[string]prowler.PluginSetting{},
	}
	// ディレクトリを再帰的に探索するための関数
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "metadata.json") {
			plugin, pluginSetting, err := extractPluginInfo(path)
			if err != nil {
				log.Printf("Failed to process file %s: %v", path, err)
			}
			setting.SpecificPluginSetting[plugin] = *pluginSetting
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	return &setting, nil
}

type pluginMetadata struct {
	Severity       string                    `json:"Severity"`
	CheckID        string                    `json:"CheckID"`
	CheckTitle     string                    `json:"CheckTitle"`
	ServiceName    string                    `json:"ServiceName"`
	SubServiceName string                    `json:"SubServiceName"`
	ResourceType   string                    `json:"ResourceType"`
	Risk           string                    `json:"Risk"`
	Remediation    pluginMetadataRemediation `json:"Remediation"`
}

type pluginMetadataRemediation struct {
	Recommendation pluginMetadataRecommendation `json:"Recommendation"`
}

type pluginMetadataRecommendation struct {
	Text string `json:"Text"`
	URL  string `json:"Url"`
}

func extractPluginInfo(filePath string) (string, *prowler.PluginSetting, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read file: %w", err)
	}

	// JSONをパース
	var metadata pluginMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	// Settingを生成
	pluginSetting := &prowler.PluginSetting{
		Tags: generateTags(metadata.ServiceName, metadata.SubServiceName, metadata.ResourceType),
		Recommend: &prowler.PluginRecommend{
			Risk:           prowler.Ptr(generateRisk(metadata.CheckTitle, metadata.Risk)),
			Recommendation: prowler.Ptr(generateRecommendation(metadata.Remediation.Recommendation.Text, metadata.Remediation.Recommendation.URL)),
		},
	}

	return fmt.Sprintf("%s/%s", metadata.ServiceName, metadata.CheckID), pluginSetting, nil
}

func generateTags(serviceName, subServiceName, resourceType string) []string {
	var tags []string
	if serviceName != "" {
		tags = append(tags, serviceName)
	}
	if subServiceName != "" {
		tags = append(tags, subServiceName)
	}
	if resourceType != "" {
		tags = append(tags, resourceType)
	}
	return tags
}

func generateRisk(title, risk string) string {
	return fmt.Sprintf("%s\n- %s", title, risk)
}

func generateRecommendation(text, url string) string {
	recommendation := text
	if url != "" {
		recommendation += fmt.Sprintf("\n- %s", url)
	}
	return recommendation
}

func updatePlugin(currentPlugin, remotePlugin *prowler.ProwlerSetting) error {
	newSetting := prowler.ProwlerSetting{
		IgnorePlugin:          currentPlugin.IgnorePlugin,
		SpecificPluginSetting: map[string]prowler.PluginSetting{}, // プラグインは空にしておく
	}

	// 削除されたプラグイン
	deletedPlugins := map[string]bool{}
	for pluginFullName := range currentPlugin.SpecificPluginSetting {
		if _, ok := remotePlugin.SpecificPluginSetting[pluginFullName]; !ok {
			deletedPlugins[pluginFullName] = true
			log.Printf("Deleted plugin: %s", pluginFullName)
		}
	}

	// プラグインをソート
	sortedPlugins := []string{}
	for pluginFullName := range remotePlugin.SpecificPluginSetting {
		sortedPlugins = append(sortedPlugins, pluginFullName)
	}
	sort.Strings(sortedPlugins)

	// プラグインを更新
	for _, pluginFullName := range sortedPlugins {
		if _, ok := deletedPlugins[pluginFullName]; ok {
			continue // 削除されたプラグインはスキップ
		}
		if _, ok := currentPlugin.SpecificPluginSetting[pluginFullName]; ok {
			// 既存のプラグインはそのまま
			current := currentPlugin.SpecificPluginSetting[pluginFullName]
			newSetting.SpecificPluginSetting[pluginFullName] = prowler.PluginSetting{
				Score:                   current.Score,
				Tags:                    current.Tags,
				SkipResourceNamePattern: current.SkipResourceNamePattern,
				IgnoreMessagePattern:    current.IgnoreMessagePattern,
				Recommend: &prowler.PluginRecommend{
					Risk:           current.Recommend.Risk,
					Recommendation: current.Recommend.Recommendation,
				},
			}
		} else {
			// 新しいプラグインの場合は追加
			new := remotePlugin.SpecificPluginSetting[pluginFullName]
			newSetting.SpecificPluginSetting[pluginFullName] = prowler.PluginSetting{
				Score:                   new.Score,
				Tags:                    new.Tags,
				SkipResourceNamePattern: new.SkipResourceNamePattern,
				IgnoreMessagePattern:    new.IgnoreMessagePattern,
				Recommend: &prowler.PluginRecommend{
					Risk:           new.Recommend.Risk,
					Recommendation: new.Recommend.Recommendation,
				},
			}
		}
	}

	// 更新されたYAMLをファイルに書き込む
	file, err := os.Create(pluginFile)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2) // インデントをスペース2つに設定
	err = encoder.Encode(newSetting)
	if err != nil {
		return fmt.Errorf("failed to encode updated YAML: %w", err)
	}
	return nil
}
