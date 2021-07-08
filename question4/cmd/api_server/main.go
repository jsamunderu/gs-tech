package main

import (
	"fmt"
	"net/http"
	"os"
	"question4/pkg/questions"

	goflags "flag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "profile-svr",
	Short: "A brief description of profile-svr",
	Long:  `A longer description that spans multiple lines...`,
	Run:   mainRun,
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.ErrorS(err, "Error starting rootCmd.Execute()")
		os.Exit(1)
	}
}

func init() {
	klog.InitFlags(nil)
	rootCmd.Flags().AddGoFlagSet(goflags.CommandLine)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.profile-svr.yaml)")

	rootCmd.PersistentFlags().String("port", "8080", "Port to listen on")
	rootCmd.PersistentFlags().String("host", "", "Host IP to listen on. If the host is empty it will listen on all IPs")
	rootCmd.PersistentFlags().String("dbHost",
		"localhost:6379", "Database storing profiles")
	rootCmd.PersistentFlags().String("profileSchema",
		"1", "Database schema storing payment messages")
	rootCmd.PersistentFlags().String("sessionSchema",
		"2", "Database schema storing sessions")
	rootCmd.PersistentFlags().String("usersSchema",
		"3", "Database schema storing users")
}

func initConfig() {
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		klog.ErrorS(err, "viper.BindPFlags")
	}

	verboseLogging := viper.GetString("v")

	viper.AutomaticEnv() // read in environment variables that match

	viper.AddConfigPath("global")
	viper.AddConfigPath("local")
	viper.AddConfigPath(".")

	viper.SetConfigName(".global-config")

	if err := viper.ReadInConfig(); err == nil {
		klog.InfoS("viper.ReadInConfig", "file", viper.ConfigFileUsed())
	} else {
		klog.ErrorS(err, "viper.ReadInConfig failed")
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".profile-svr")
	}

	if err := viper.MergeInConfig(); err == nil {
		klog.InfoS("viper.MergeInConfig", "file", viper.ConfigFileUsed())
	} else {
		klog.ErrorS(err, "viper.MergeInConfig failed")
	}

	if verboseLogging == "0" {
		if err := goflags.Set("v", viper.GetString("v")); err != nil {
			klog.ErrorS(err, "goflags.Set")
		}
	}
	for _, v := range viper.AllKeys() {
		klog.InfoS("Configs loaded", v, viper.Get(v))
	}
}

func main() {
	execute()
}

func mainRun(cmd *cobra.Command, args []string) {
	defer klog.Flush()

	profile, err := questions.NewProfile(viper.GetString("dbHost"),
		viper.GetString("profilesSchema"), viper.GetString("sessionsSchema"), viper.GetString("usersSchema"))
	if err != nil {
		klog.ErrorS(err, "Setup error")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", profile.DefaultPath)
	mux.HandleFunc("/login", profile.Login)
	mux.HandleFunc("/profileUpdate", profile.ProfileUpdate)
	mux.HandleFunc("/profileDelete", profile.ProfileDelete)
	mux.HandleFunc("/query", profile.GetProfile)

	svr := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "localhost", "8080"),
		Handler: mux,
	}

	if err := svr.ListenAndServe(); err != nil {
		fmt.Println(err, "Error!")
	}
}
