package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rshulabs/HgCache/pkg/errorx"
	"github.com/rshulabs/HgCache/pkg/logx"
	"github.com/rshulabs/HgCache/pkg/utils/term"
	"github.com/rshulabs/HgCache/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	v *bool
)

type App struct {
	basename    string // 打包名称
	name        string // 应用名称
	description string
	options     CliOptions
	runFunc     RunFunc
	noVersion   bool
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs // 位置参数校验
	cmd         *cobra.Command
}

type RunFunc func(basename string) error

type Option func(app *App)

func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func NewApp(name string, basename string, opts ...Option) *App {
	app := &App{
		name:     name,
		basename: basename,
	}
	for _, opt := range opts {
		opt(app)
	}
	app.buildCommand()
	return app
}

// 命令行参数设置
func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// 出现错误 停止打印usage
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	// 参数选项
	cmd.Flags().SortFlags = true
	// 子命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}
	// 项目运行函数
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}
	// 命令行参数选项
	var flagSets FlagSets
	if a.options != nil {
		flagSets = a.options.Flags() // 绑定所有命令行参数
		fs := cmd.Flags()
		for _, f := range flagSets.Sets {
			fs.AddFlagSet(f)
		}
	}
	// TODO 版本
	if !a.noVersion {
		v = version.AddVersionFlag(cmd.Name(), flagSets.FlagSet("global"))
	}
	// 配置文件
	if !a.noConfig {
		// 读取配置文件
		addConfigFlag(flagSets.FlagSet("global"))
	}
	// help
	addHelpFlag(cmd.Name(), flagSets.FlagSet("global"))
	cmd.Flags().AddFlagSet(flagSets.FlagSet("global"))
	// 打印输出
	addCmdTemplate(&cmd, flagSets)
	a.cmd = &cmd
}

// 执行命令
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	// 版本
	if *v {
		version.FullVersion()
	}
	printWorkingDir()
	// 配置
	if !a.noConfig {
		// 绑定所有命令行参数
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	// 命令行参数选项
	if a.options != nil {
		// 验证
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// 运行函数
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	logx.Infof("WorkingDir: %s", wd)
}

func (a *App) applyOptionRules() error {
	// if completeableOptions, ok := a.options.(CompleteableOptions); ok {
	// 	if err := completeableOptions.Complete(); err != nil {
	// 		return err
	// 	}
	// }

	if errs := a.options.Validate(); len(errs) != 0 {
		return errorx.NewAggregate(errs)
	}

	return nil
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets FlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
