package util

import (
	"fmt"
	"github.com/coreos/etcd/pkg/testutil"
	"testing"
)

func TestRemoveElements(t *testing.T) {
	arr := []string{"1", "2", "3", "4"}
	arri := make([]interface{}, len(arr))
	for i, e := range arr {
		arri[i] = e
	}
	elements := RemoveElements(arri, "2", "3")
	fmt.Println(elements)
	testutil.AssertEqual(t, "1", elements[0])
	testutil.AssertEqual(t, "4", elements[1])
}
