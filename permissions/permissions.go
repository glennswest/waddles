package permissions

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

var (
	permSystem PermissionSystem
)

type permissionSet struct {
	Nodes       []*permissionNode `toml:"nodes"`
	Groups      []*permissionGroup
	Name        string
	Description string
}

type permissionGroup struct {
	Sets        []*permissionSet `toml:"sets"`
	Name        string
	Description string
	RoleID      string
}

type permissionTree struct {
	Sets   []permissionSet
	Groups []permissionGroup
}

type permissionNode struct {
	Sets       []*permissionSet
	Identifier string
}

func (group *permissionGroup) UnmarshalTOML(i interface{}) error {
	iMap, ok := i.(map[string]interface{})

	if !ok {
		return fmt.Errorf("type assertion error: wants %T, have %T", map[string]interface{}{}, i)
	}

	name, _ := iMap["name"]
	group.Name = name.(string)

	description, _ := iMap["description"]
	group.Description = description.(string)

	role, _ := iMap["role"]
	group.RoleID = role.(string)

	rawSets, _ := iMap["sets"].([]interface{})

	for _, rawSet := range rawSets {
		pm := CurrentPermissionSystem()
		setString, _ := rawSet.(string)
		set := pm.GetSetFromName(setString)

		set.Groups = append(set.Groups, group)
		group.Sets = append(group.Sets, set)
	}

	return nil
}

func (set *permissionSet) UnmarshalTOML(i interface{}) error {
	iMap, ok := i.(map[string]interface{})

	if !ok {
		return fmt.Errorf("type assertion error: wants %T, have %T", map[string]interface{}{}, i)
	}

	name, _ := iMap["name"]
	set.Name = name.(string)

	description, _ := iMap["description"]
	set.Description = description.(string)

	rawNodes, _ := iMap["nodes"].([]interface{})
	pm := CurrentPermissionSystem()

	for _, rawNode := range rawNodes {
		nodeString, _ := rawNode.(string)
		node := pm.GetNodeFromIdentifier(nodeString)

		node.Sets = append(node.Sets, set)
		set.Nodes = append(set.Nodes, node)
	}

	pm.Tree.Sets = append(pm.Tree.Sets, *set)

	return nil
}

func parsePermissionConfig(tomlBytes []byte) (permissionTree, error) {
	permTree := permissionTree{}

	err := toml.Unmarshal(tomlBytes, &permTree)

	return permTree, err
}