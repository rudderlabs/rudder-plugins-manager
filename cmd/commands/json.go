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
	Plugin string `json:"plugin"`
	Data   any    `json:"data"`
	// Transformation function is applied based on metadata
	MetaData any `json:"metadata"`
}

type jsonTransformer struct {
	pluginsManager *plugins.Manager
}

func (c *jsonTransformer) transformRecords(records []Record, provider string) ([]Record, error) {
	for i, record := range records {
		plugin, err := c.pluginsManager.GetPlugin(provider, record.Plugin)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		metadata := record.MetaData
		if metadata == nil {
			metadata = record.Data
		}
		transformer, err := plugin.GetTransformer(metadata)
		if err != nil {
			fmt.Printf("get transformer with data %v of plugin %s was failed with error %v\n", metadata, record.Plugin, err)
			continue
		}
		records[i].Data, err = transformer.Transform(record.Data)
		if err != nil {
			fmt.Printf("transforming the data %v for plugin %s was failed with error %v\n", record.Data, record.Plugin, err)
			continue
		}
	}
	return records, nil
}

func (c *jsonTransformer) transform(cmd *cobra.Command, args []string) error {
	records, err := utils.ReadRecordsFromJSONFile[Record](cmd.Flag("input").Value.String())
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
	if err = os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		return err
	}
	err = os.WriteFile(outputFile, output, 0o644)
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
	if err := jsonCobraCmd.MarkFlagRequired("input"); err != nil {
		panic(err)
	}
	jsonCobraCmd.Flags().StringP("output", "o", "output.json", "output json file")
	jsonCobraCmd.Flags().StringP("provider", "p", "", "provider name")
	if err := jsonCobraCmd.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}
	return jsonCobraCmd
}
