package nvelope_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/muir/nvelope"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.Equal(t, http.StatusNotModified, nvelope.GetReturnCode(nvelope.ReturnCode(fmt.Errorf("x"), http.StatusNotModified)), "unwrapped")
	assert.Equal(t, http.StatusSeeOther, nvelope.GetReturnCode(errors.Wrap(nvelope.ReturnCode(fmt.Errorf("x"), http.StatusSeeOther), "o")), "wrapped")
	assert.Equal(t, http.StatusBadRequest, nvelope.GetReturnCode(nvelope.BadRequest(fmt.Errorf("x"))), "bad")
	assert.Equal(t, http.StatusUnauthorized, nvelope.GetReturnCode(nvelope.Unauthorized(fmt.Errorf("x"))), "unauth")
	assert.Equal(t, http.StatusForbidden, nvelope.GetReturnCode(nvelope.Forbidden(fmt.Errorf("x"))), "forbid")
}
