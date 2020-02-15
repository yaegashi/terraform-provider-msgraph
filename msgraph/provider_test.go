package msgraph_test

import (
	"testing"

	"github.com/yaegashi/terraform-provider-msgraph/msgraph"
)

func TestProvider(t *testing.T) {
	if err := msgraph.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
