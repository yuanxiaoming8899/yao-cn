package sui

func init() {
	WatchCmd.PersistentFlags().StringVarP(&data, "data", "d", "::{}", L("Session Data"))
}
