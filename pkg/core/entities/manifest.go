package entities

import "github.com/mrparano1d/noxite/pkg/core/fields"

type PackageVersion struct {
	Name                 fields.RequiredString
	Version              fields.RequiredString
	Description          *string
	Keywords             []fields.RequiredString
	Homepage             *fields.Website
	Bugs                 *fields.Bugs
	License              *string
	Author               *fields.MixedAuthor
	Contributors         fields.MixedAuthors
	Funding              []fields.UrlType
	Files                []fields.RequiredString
	Main                 *fields.RequiredString
	Browser              *fields.RequiredString
	Bin                  map[fields.RequiredString]fields.RequiredString
	Man                  []fields.RequiredString
	Directories          *fields.Directories
	Repository           *fields.Repository
	Scripts              map[fields.RequiredString]fields.RequiredString
	Config               map[fields.RequiredString]fields.RequiredString
	Dependencies         map[fields.RequiredString]fields.RequiredString
	DevDependencies      map[fields.RequiredString]fields.RequiredString
	PeerDependencies     map[fields.RequiredString]fields.RequiredString
	PeerDependenciesMeta map[fields.RequiredString]map[fields.RequiredString]any
	BundledDependencies  []fields.RequiredString
	OptionalDependencies map[fields.RequiredString]fields.RequiredString
	Overrides            map[fields.RequiredString]fields.RequiredString
	Engines              map[fields.RequiredString]fields.RequiredString
	OS                   []fields.RequiredString
	CPU                  []fields.RequiredString
	Private              *bool
	PublishConfig        map[fields.RequiredString]any
	Workspaces           []fields.RequiredString

	Integrity   fields.RequiredString
	SHASUM      fields.RequiredString
	ContentType fields.RequiredString
	Data        fields.RequiredString
	Length      int
	Readme      *string
}
