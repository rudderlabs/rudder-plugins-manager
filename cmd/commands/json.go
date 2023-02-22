package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rudderlabs/rudder-transformations/plugins"
	"github.com/rudderlabs/rudder-transformations/utils"
	"github.com/spf13/cobra"
)

type Record struct {
	Plugin string      `json:"plugin"`
	Data   interface{} `json:"data"`
	// Transformation function is applied based on metadata
	MetaData interface{} `json:"metadata"`
}

type jsonTransformer struct {
	pluginsManager *plugins.Manager
}

func (c *jsonTransformer) transformRecords(records []Record, provider string) ([]Record, error) {
	for i, record := range records {
		plugin, ok := c.pluginsManager.GetPlugin(provider, record.Plugin)
		if !ok {
			fmt.Printf("plugin %s not found\n", record.Plugin)
			continue
		}
		metadata := record.MetaData
		if metadata == nil {
			metadata = record.Data
		}
		transformer, err := plugin.GetTransformer(metadata)
		if err != nil {
			fmt.Printf("failed to get transformer with data %v of plugin %s\n", metadata, record.Plugin)
			continue
		}
		records[i].Data, err = transformer.Transform(record.Data)
		if err != nil {
			fmt.Printf("failed to transform the data %v for plugin %s\n", record.Data, record.Plugin)
			continue
		}
	}
	return records, nil
}

func (c *jsonTransformer) transform(cmd *cobra.Command, args []string) error {
	records, err := utils.ReadJSONFromFile[Record](cmd.Flag("input").Value.String())
	if err != nil {
		return err
	}
	provider := cmd.Flag("provider").Value.String()
	records, err = c.transformRecords(records, provider)
	if err != nil {
		return err
	}
	output, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	outputFile := cmd.Flag("output").Value.String()
	if err = os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}
	err = os.WriteFile(outputFile, output, 0644)
	return err
}

func GetJSONCmd(manager *plugins.Manager) *cobra.Command {
	jsonTransformer := &jsonTransformer{pluginsManager: manager}
	jsonCobraCmd := &cobra.Command{
		Use:   "json",
		Short: "Transforms json file",
		Long:  `Transforms json file using the plugin from input file and provider flag.`,
		RunE:  jsonTransformer.transform,
	}
	jsonCobraCmd.Flags().StringP("input", "i", "", "input json file")
	jsonCobraCmd.MarkFlagRequired("input")
	jsonCobraCmd.Flags().StringP("output", "o", "output.json", "output json file")
	jsonCobraCmd.Flags().StringP("provider", "p", "", "provider name")
	jsonCobraCmd.MarkFlagRequired("provider")
	return jsonCobraCmd
}
