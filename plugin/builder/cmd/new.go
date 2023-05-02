/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "创建一个新的插件",
	Long:  `创建一个新的插件`,
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := promptui.Prompt{
			Label:       "plugin name",
			HideEntered: true,
			Validate:    ValidAllAlpha,
		}
		var err error
		pluginInfo := &BuildMetaInfo{}
		pluginInfo.PluginName, err = prompt.Run()
		if err != nil {
			return err
		}
		prompt = promptui.Prompt{Label: "plugin description", HideEntered: true}
		pluginInfo.Description, err = prompt.Run()
		if err != nil {
			return err
		}
		prompt = promptui.Prompt{Label: "author", HideEntered: true}
		pluginInfo.Author, err = prompt.Run()
		if err != nil {
			return err
		}
		prompt = promptui.Prompt{Label: "author url", HideEntered: true}
		pluginInfo.Url, err = prompt.Run()
		if err != nil {
			return err
		}
		err = os.MkdirAll(pluginInfo.PluginName, 0777)
		if err != nil {
			return err
		}
		goFile = strings.ReplaceAll(goFile, "{{.pluginName}}", pluginInfo.PluginName)
		err = os.WriteFile(filepath.Join(pluginInfo.PluginName, pluginInfo.PluginName+".go"), []byte(goFile), 0777)
		if err != nil {
			return err
		}
		GoCmd := exec.Command("go", "mod", "init", pluginInfo.PluginName)
		GoCmd.Dir = "./" + pluginInfo.PluginName
		GoCmd.Stderr = os.Stderr
		err = GoCmd.Run()
		if err != nil {
			return err
		}
		GoCmd = exec.Command("go", "mod", "tidy")
		GoCmd.Dir = "./" + pluginInfo.PluginName
		GoCmd.Stderr = os.Stderr
		err = GoCmd.Run()
		if err != nil {
			return err
		}
		metaInfo, err := json.Marshal(pluginInfo)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = json.Indent(&buf, metaInfo, "", "\t")
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(pluginInfo.PluginName, "meta.json"), buf.Bytes(), 0777)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ValidAllAlpha(s string) error {
	if match, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, s); match {
		return nil
	}
	return fmt.Errorf("err")
}

var goFile = `//go:build tinygo.wasm

package main

import (
	"github.com/opq-osc/Yui/proto"
	"context"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
)

type {{.pluginName}} struct {
}

func (p {{.pluginName}}) OnRemoteCallEvent(ctx context.Context, req *proto.RemoteCallReq) (*proto.RemoteCallReply, error) {
	//TODO implement me
	panic("implement me")
}

func (p {{.pluginName}}) OnCronEvent(ctx context.Context, req *proto.CronEventReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p {{.pluginName}}) OnFriendMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p {{.pluginName}}) OnPrivateMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p {{.pluginName}}) OnGroupMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (p {{.pluginName}}) Unload(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil

}
func (p {{.pluginName}}) Init(ctx context.Context, _ *emptypb.Empty) (*proto.InitReply, error) {
	return &proto.InitReply{
		Ok:      true,
		Message: "Success",
	}, nil
}

func main() {
	proto.RegisterEvent({{.pluginName}}{})
}
`
