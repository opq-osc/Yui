/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	tinyGoPath string
	output     string
	privateKey string
)

type BuildMetaInfo struct {
	meta.PluginMeta
	Permissions []string `json:"Permissions"`
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build 插件的go文件路径",
	Short: "创建OPQ插件文件",
	Long:  `创建OPQ插件文件`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取 meta 信息
		path := filepath.Dir(args[0])
		metaByte, err := os.ReadFile(filepath.Join(path, "meta.json"))
		if err != nil {
			return err
		}
		var rawM = BuildMetaInfo{}
		err = json.Unmarshal(metaByte, &rawM)
		if err != nil {
			return err
		}
		// 权限转换
		var m = meta.PluginMeta{
			PluginName:   rawM.PluginName,
			Description:  rawM.Description,
			Dependencies: rawM.Dependencies,
			Author:       rawM.Author,
			Url:          rawM.Url,
			Version:      rawM.Version,
			Permissions:  meta.SumPermissions(rawM.Permissions...),
			Sha256:       rawM.Sha256,
			Sign:         false,
			SignInfo:     rawM.SignInfo,
		}

		if output == "" {
			output = filepath.Join(path, m.PluginName)
		}
		err = RunCmd(tinyGoPath, "build", "-o", output+".wasm", "-scheduler=none", "-target=wasi", "--no-debug", args[0])
		if err != nil {
			return err
		}
		defer os.Remove(output + ".wasm")
		b, err := os.ReadFile(output + ".wasm")
		if err != nil {
			return err
		}
		s := sha256.New()
		s.Write(b)
		sha256hash := hex.EncodeToString(s.Sum(nil))
		m.Sha256 = sha256hash

		var headerBuf bytes.Buffer
		e := gob.NewEncoder(&headerBuf)
		if err := e.Encode(m); err != nil {
			panic(err)
		}
		header := headerBuf.Bytes()
		var buf bytes.Buffer
		buf.Write([]byte("OPQ"))
		if err := binary.Write(&buf, binary.LittleEndian, int32(meta.PluginApiVersion)); err != nil {
			panic(err)
		}
		if privateKey == "" {
			if err := binary.Write(&buf, binary.LittleEndian, int32(0)); err != nil {
				panic(err)
			}
		} else {
			if err := binary.Write(&buf, binary.LittleEndian, int32(1)); err != nil {
				panic(err)
			}
			// 签名
			privateKeyByte, err := hex.DecodeString(privateKey)
			if err != nil {
				panic(err)
			}
			key, err := crypto.ToECDSA(privateKeyByte)
			if err != nil {
				panic(err)
			}
			sha := sha256.New()
			sha.Write(headerBuf.Bytes())
			hash := sha.Sum(nil)
			r, s, err := ecdsa.Sign(rand.Reader, key, hash)
			if err != nil {
				panic(err)
			}
			Rbyte, err := r.MarshalText()
			if err != nil {
				panic(err)
			}
			Sbyte, err := s.MarshalText()
			if err != nil {
				panic(err)
			}
			if err := binary.Write(&buf, binary.LittleEndian, int32(len(Rbyte))); err != nil {
				panic(err)
			}
			buf.Write(Rbyte)
			if err := binary.Write(&buf, binary.LittleEndian, int32(len(Sbyte))); err != nil {
				panic(err)
			}
			buf.Write(Sbyte)
		}
		if err := binary.Write(&buf, binary.LittleEndian, int32(len(header))); err != nil {
			panic(err)
		}
		buf.Write(header)
		buf.Write(b)
		f, err := os.OpenFile(output+".opq", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = io.Copy(f, &buf)
		if err != nil {
			panic(err)
		}
		return nil
	},
}

func RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir, _ = os.Getwd()
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	viper.AutomaticEnv()
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&tinyGoPath, "tinyGoPath", "t", "tinygo", "设置tinyGo程序路径")
	buildCmd.Flags().StringVarP(&output, "output", "o", "", "设置输出插件文件名，不包含后缀")
	buildCmd.Flags().StringVarP(&privateKey, "privateKey", "p", viper.GetString("OPQ_KEY"), "为插件签名")
}
