package lb

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	hash := NewMap(3, nil)
	hash.Add("192.168.60.113", "192.168.60.5", "192.168.60.65")
	testCliIps := []string{"192.168.23.5", "192.168.3.6", "192.168.34.67"}
	for _, v := range testCliIps {
		fmt.Printf("the node of  [%s] is [%s].\n",v,hash.Get(v))
	}
}
