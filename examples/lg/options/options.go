package options

import "github.com/brezzgg/go-packages/lg"

func ExampleLevelOptions() {
	lg.Log(lg.NewLogLevel(lg.ClrFgYellow, "!Example!").WithOptions(
		lg.LevelOptionDisableCaller, lg.LevelOptionDisableTime),
		"Log level options example",
	)

	level := lg.NewLogLevel(lg.ClrFgBlue, "LogLevelOption")

	lg.Log(level.WithOptions(), "some text", lg.C{"option": ""})

	lg.Log(level.WithOptions(lg.LevelOptionDisableTime), "some text", lg.C{"option": "LevelOptionDisableTime"})
	lg.Log(level.WithOptions(lg.LevelOptionTimeDisableOffset), "some text", lg.C{"option": "LevelOptionTimeDisableOffset"})

	lg.Log(level.WithOptions(lg.LevelOptionDisableCaller), "some text", lg.C{"option": "LevelOptionDisableCaller"})
	lg.Log(level.WithOptions(lg.LevelOptionCallerOnlyFile), "some text", lg.C{"option": "LevelOptionCallerOnlyFile"})
	lg.Log(level.WithOptions(lg.LevelOptionCallerOnlyFunc), "some text", lg.C{"option": "LevelOptionCallerOnlyFunc"})
	lg.Log(level.WithOptions(lg.LevelOptionCallerDisableLine), "some text", lg.C{"option": "LevelOptionCallerDisableLine"})
	lg.Log(level.WithOptions(lg.LevelOptionCallerDisableFunc), "some text", lg.C{"option": "LevelOptionCallerDisableFunc"})
	lg.Log(level.WithOptions(lg.LevelOptionCallerDisableFile), "some text", lg.C{"option": "LevelOptionCallerDisableFile"})

	lg.Log(level.WithOptions(lg.LevelOptionDisableTime, lg.LevelOptionDisableCaller), "some text", lg.C{"option": "LevelOptionDisableTime,LevelOptionDisableCaller"})
	lg.Log(level.WithOptions(lg.LevelOptionTimeDisableOffset), "some text", lg.C{"option": "LevelOptionTimeDisableOffset"})
}
