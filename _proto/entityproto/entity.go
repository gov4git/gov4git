package entityproto

import (
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/proto/govproto"
)

func EntityDirpath(namespace string, entity string) string {
	return filepath.Join(govproto.GovRoot, namespace, entity)
}

const ValueFilebase = "value.json"

func ValueFilepath(namespace string, entity string) string {
	return filepath.Join(EntityDirpath(namespace, entity), ValueFilebase)
}

const PropDirbase = "prop"

func PropValueFilepath(namespace string, entity string, property string) string {
	return filepath.Join(EntityDirpath(namespace, entity), PropDirbase, property)
}

func BalancePropKey(balance string) string {
	return fmt.Sprintf("balance:%v", balance)
}
