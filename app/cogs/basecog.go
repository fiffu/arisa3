package cogs

// func (c *DefaultCog) registerCommands(ctx context.Context, s *dgo.Session) error {
// 	err := c.registry.Register(
// 		ctx,
// 		s,
// 		types.NewCommand("roll").ForChat().
// 			Desc("Rolls dice (supports algebraic notation, such as !roll 3d5+10)").
// 			Options(types.NewOption("expression").String()).
// 			Handler(c.roll),
// 	)
// 	c.registry.Finalise(ctx, s)
// 	return err
// }

// func (c *DefaultCog) roll(s *dgo.Session, i *dgo.InteractionCreate, cmd types.ICommand, args types.IArgs) error {
// 	ctx := c.util.CommandContext(cmd, i.User)
// 	resp := types.NewResponse()
// 	if expr, ok := args.String("expression"); ok {
// 		c.app.Infof(ctx, "expression='%s'", expr)
// 		lo, hi, num := roll(expr)
// 		resp.Content(fmt.Sprintf("Rolling %d-%d: %d", lo, hi, num))
// 	} else {
// 		c.app.Infof(ctx, "missing expression")
// 		resp.Content("missing expression")
// 	}
// 	return s.InteractionRespond(i.Interaction, resp.Data())
// }

// func roll(expr string) (int, int, int) {
// 	lo, hi := 0, 99
// 	return lo, hi, rollRand(lo, hi)
// }

// func rollRand(lo, hi int) int {
// 	return rand.Intn(lo+hi) - lo
// }
