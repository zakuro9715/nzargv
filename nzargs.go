package nzargs

import (
	"fmt"
	"os"
	"strings"
)

type App struct {
	FlagsValueN map[string]int
}

func New() *App {
	return &App{map[string]int{}}
}

func (app *App) FlagN(name string, n int) *App {
	app.FlagsValueN[name] = n
	return app
}

func (app *App) GetFlagN(name string) int {
	n, ok := app.FlagsValueN[name]
	if ok {
		return n
	}
	return 0
}

func (app *App) NormalizeArgs() {
	app.Normalize(os.Args[1:])
}

func checkFlagValue(name string, n int, args []string) error {
	i := 0
	for ; i < len(args) && !strings.HasPrefix(args[i], "-"); i++ {
	}
	if i < n {
		format := "Flag %v values is too few. %v values required but has %v args."
		return fmt.Errorf(format, name, n, i)
	}
	return nil
}

func (app *App) Normalize(argv []string) ([]Value, error) {
	normalized := make([]Value, 0)
	for i := 0; i < len(argv); i++ {
		v := argv[i]
		switch {
		case strings.HasPrefix(v, "--"):
			v = strings.TrimPrefix(v, "--")
			splited := strings.SplitN(v, "=", 2)[0:2]
			name := splited[0]
			if len(splited[1]) == 0 {
				n := app.GetFlagN(name)
				if err := checkFlagValue(name, n, argv[i+1:]); err != nil {
					return nil, err
				}
				normalized = append(normalized, NewFlag(name, argv[i+1:i+1+n]...))
				i += n
			} else {
				values := strings.Split(splited[1], ",")
				normalized = append(normalized, NewFlag(name, values...))
			}
		case strings.HasPrefix(v, "-"):
			v = strings.TrimPrefix(v, "-")
			splited := strings.SplitN(v, "=", 2)[0:2]
			names := splited[0]
			last_name := string(names[len(names)-1])
			for _, name := range names[:len(names)-1] {
				normalized = append(normalized, NewFlag(string(name)))
			}
			if len(splited[1]) == 0 {
				n := app.GetFlagN(last_name)
				if err := checkFlagValue(last_name, n, argv[i+1:]); err != nil {
					return nil, err
				}
				normalized = append(normalized, NewFlag(last_name, argv[i+1:i+1+n]...))
				i += n
			} else {
				values := strings.Split(splited[1], ",")
				normalized = append(normalized, NewFlag(last_name, values...))
			}
		default:
			normalized = append(normalized, NewArg(v))
		}
	}
	return normalized, nil
}

func (app *App) NormalizeToString(argv []string) ([]string, error) {
	normalized, err := app.Normalize(argv)
	if err != nil {
		return nil, err
	}
	values := make([]string, len(normalized))
	for i, v := range normalized {
		values[i] = v.Text()
	}
	return values, nil
}
