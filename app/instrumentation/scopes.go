package instrumentation

import "github.com/fiffu/arisa3/lib"

// supportedScope
type supportedScope string

// Supported trace scopes
const (
	commandScope      supportedScope = "arisa3/command"
	eventScope        supportedScope = "arisa3/event"
	databaseScope     supportedScope = "database"
	externalHTTPScope supportedScope = "external-http"
	vendorScope       supportedScope = "vendor"
)

type Command string

func (sn Command) scope() supportedScope { return commandScope }
func (sn Command) name() string          { return string(sn) }

type Event string

func (sn Event) scope() supportedScope { return eventScope }
func (sn Event) name() string          { return string(sn) }

type Database string

func (sn Database) scope() supportedScope { return databaseScope }
func (sn Database) name() string          { return string(sn) }

type ExternalHTTP string

func (sn ExternalHTTP) scope() supportedScope { return externalHTTPScope }
func (sn ExternalHTTP) name() string          { return string(sn) }

func Vendor(callable any) ScopedName { return vendor(lib.FuncName(callable)) }

type vendor string

func (sn vendor) scope() supportedScope { return vendorScope }
func (sn vendor) name() string          { return string(sn) }
