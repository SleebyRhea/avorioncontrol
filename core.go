package main

// Core only exists for logging purposes, and contains no other state
type Core struct {
	loglevel int
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of an alliance
func (c *Core) UUID() string {
	return "Core"
}

// Loglevel returns the loglevel of an alliance
func (c *Core) Loglevel() int {
	return config.Loglevel()
}

// SetLoglevel sets the loglevel of an alliance
func (c *Core) SetLoglevel(l int) {
	config.SetLoglevel(l)
}
