package ntlm

import (
	"io/ioutil"
	"os"
	"testing"

	h "github.com/koltyakov/gosip-sandbox/test/helpers"
	u "github.com/koltyakov/gosip-sandbox/test/utils"
)

var cnfgPath = "./config/private.ntlm.json"

func TestGettingAuth(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No config found, skipping...")
	}
	err := h.CheckAuth(
		&AuthCnfg{},
		cnfgPath,
		[]string{"SiteURL", "Username", "Password", "Domain"},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestBasicRequest(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckRequest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}

func TestGettingDigest(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthEdgeCases(t *testing.T) {

	t.Run("ReadConfig/MissedConfig", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		if err := cnfg.ReadConfig("wrong_path.json"); err == nil {
			t.Error("wrong_path config should not pass")
		}
	})

	t.Run("ReadConfig/MissedConfig", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		folderPath := u.ResolveCnfgPath("./tmp")
		filePath := u.ResolveCnfgPath("./tmp/private.ntlm.malformed.json")
		_ = os.MkdirAll(folderPath, os.ModePerm)
		_ = ioutil.WriteFile(filePath, []byte("not a json"), 0644)
		if err := cnfg.ReadConfig(filePath); err == nil {
			t.Error("malformed config should not pass")
		}
		_ = os.RemoveAll(filePath)
	})

	t.Run("WriteConfig", func(t *testing.T) {
		folderPath := u.ResolveCnfgPath("./tmp")
		filePath := u.ResolveCnfgPath("./tmp/private.ntlm.json")
		cnfg := &AuthCnfg{SiteURL: "test"}
		_ = os.MkdirAll(folderPath, os.ModePerm)
		if err := cnfg.WriteConfig(filePath); err != nil {
			t.Error(err)
		}
		_ = os.RemoveAll(filePath)
	})

	t.Run("SetMasterkey", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		cnfg.SetMasterkey("key")
		if cnfg.masterKey != "key" {
			t.Error("unable to set master key")
		}
	})

	t.Run("GetAuth", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		r, err := cnfg.GetAuth()
		if err != nil {
			t.Error(err)
		}
		if r != "" {
			t.Error("ntml's t.GetAuth should not return anything")
		}
	})

}
