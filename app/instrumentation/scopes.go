package instrumentation

import (
	"fmt"

	"github.com/fiffu/arisa3/lib"
)

// supportedScope
type supportedScope string

// Supported trace scopes
const (
	internalScope     supportedScope = "arisa3/internal"
	commandScope      supportedScope = "arisa3/command"
	eventScope        supportedScope = "arisa3/event"
	databaseScope     supportedScope = "database"
	externalHTTPScope supportedScope = "external-http"
	vendorScope       supportedScope = "vendor"
)

type Internal string

func (sn Internal) scope() supportedScope { return internalScope }
func (sn Internal) name() string          { return fmt.Sprintf("Internal: %s", sn) }

type Command string

func (sn Command) scope() supportedScope { return commandScope }
func (sn Command) name() string          { return fmt.Sprintf("Command: /%s", sn) }

func EventHandler(callable any) ScopedName { return event(lib.FuncName(callable)) }

type event string

func (sn event) scope() supportedScope { return eventScope }
func (sn event) name() string          { return string(sn) }

type Database string

func (sn Database) scope() supportedScope { return databaseScope }
func (sn Database) name() string          { return fmt.Sprintf("DB: %s", sn) }

type ExternalHTTP string

func (sn ExternalHTTP) scope() supportedScope { return externalHTTPScope }
func (sn ExternalHTTP) name() string          { return string(sn) }

func Vendor(callable any) ScopedName { return vendor(lib.FuncName(callable)) }

type vendor string

func (sn vendor) scope() supportedScope { return vendorScope }
func (sn vendor) name() string          { return string(sn) }
