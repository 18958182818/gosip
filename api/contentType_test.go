package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var ctID = ""

func TestContentType(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()

	t.Run("Conf", func(t *testing.T) {
		ct := web.ContentTypes().GetByID("")
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			g := ct.Conf(preset)
			if g.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		ct, err := getRandomCT()
		if err != nil {
			t.Error(err)
		}
		if _, err := ct.Select("*,Fields/*").Expand("Fields").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		resp, err := web.ContentTypes().Top(5).Get()
		if err != nil {
			t.Error(err)
		}
		cts := resp.Data()
		if len(cts) != 5 {
			t.Error("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			t.Error("can't get content type info")
		}
		if _, err := web.ContentTypes().GetByID(cts[0].Data().ID).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("UpdateDelete", func(t *testing.T) {
		guid := uuid.New().String()
		ctID := "0x0100" + strings.ToUpper(strings.Replace(guid, "-", "", -1))
		ct := []byte(`{
			"Description":"",
			"Group":"Custom Content Types",
			"Id":{"StringValue":"` + ctID + `"},
			"Name":"test-temp-ct ` + guid + `"
		}`)
		if _, err := web.ContentTypes().Add(ct); err != nil {
			t.Error(err)
		}
		if _, err := web.ContentTypes().GetByID(ctID).Update([]byte(`{"Description":"Test"}`)); err != nil {
			t.Error(err)
		}
		if err := web.ContentTypes().GetByID(ctID).Delete(); err != nil {
			t.Error(err)
		}
	})

	// ToDo:
	// Recycle

}

func getRandomCT() (*ContentType, error) {
	sp := NewSP(spClient)
	if ctID == "" {
		resp, err := sp.Web().ContentTypes().Top(1).Get()
		if err != nil {
			return nil, err
		}
		cts := resp.Data()
		if len(cts) != 1 {
			return nil, fmt.Errorf("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			return nil, fmt.Errorf("can't get content type info")
		}
		ctID = cts[0].Data().ID
	}
	return sp.Web().ContentTypes().GetByID(ctID), nil
}
