package mongoStorage

import (
	"github.com/EdmundMartin/boselecta/pkg/flag"
	"go.mongodb.org/mongo-driver/bson"
)

func decodeBSON(result bson.M) *flag.FeatureFlag {
	fl := flag.FeatureFlag{}
	fl.Namespace = result["Namespace"].(string)
	fl.FlagName = result["FlagName"].(string)
	fl.Value = result["Value"]
	fl.Refresh = int(result["Refresh"].(int32))
	fl.Type = flag.ToFlagType(result["Type"].(string))
	return &fl
}