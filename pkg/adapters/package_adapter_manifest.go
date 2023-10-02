package adapters

import (
	"fmt"
	"strings"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

type bugs struct {
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

type contributor struct {
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
	URL   *string `json:"url,omitempty"`
}

type directories struct {
	Man string `json:"man"`
	Bin string `json:"bin"`
}

type repository struct {
	URL       string  `json:"url"`
	Type      *string `json:"type,omitempty"`
	Directory *string `json:"directory,omitempty"`
}

type dist struct {
	Tarball   string `json:"tarball"`
	Integrity string `json:"integrity"`
	SHASUM    string `json:"shasum"`
}

type attachment struct {
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
	Length      int    `json:"length"`
}

type revision struct {
	Name                 string                    `json:"name"`
	Version              string                    `json:"version"`
	Description          string                    `json:"description,omitempty"`
	Keywords             []string                  `json:"keywords,omitempty"`
	Homepage             string                    `json:"homepage,omitempty"`
	Bugs                 *bugs                     `json:"bugs,omitempty"`
	License              string                    `json:"license,omitempty"`
	Author               any                       `json:"author,omitempty"`
	Contributors         []contributor             `json:"contributors,omitempty"`
	Funding              any                       `json:"funding"`
	Files                []string                  `json:"files"`
	Main                 string                    `json:"main"`
	Browser              string                    `json:"browser"`
	Bin                  map[string]string         `json:"bin"`
	Man                  []string                  `json:"man"`
	Directories          *directories              `json:"directories"`
	Repository           *repository               `json:"repository"`
	Scripts              map[string]string         `json:"scripts"`
	Config               map[string]string         `json:"config"`
	Dependencies         map[string]string         `json:"dependencies"`
	DevDependencies      map[string]string         `json:"devDependencies"`
	PeerDependencies     map[string]string         `json:"peerDependencies"`
	PeerDependenciesMeta map[string]map[string]any `json:"peerDependenciesMeta"`
	BundledDependencies  []string                  `json:"bundledDependencies"`
	OptionalDependencies map[string]string         `json:"optionalDependencies"`
	Engines              map[string]string         `json:"engines"`
	OS                   []string                  `json:"os"`
	CPU                  []string                  `json:"cpu"`
	Private              *bool                     `json:"private,omitempty"`
	PublishConfig        map[string]any            `json:"publishConfig"`
	Workspaces           []string                  `json:"workspaces"`
	Readme               string                    `json:"readme"`
	Dist                 dist                      `json:"dist"`
}

type manifest struct {
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Readme      string                `json:"readme,omitempty"`
	Versions    map[string]revision   `json:"versions"`
	Attachments map[string]attachment `json:"_attachments"`
	DistTags    map[string]string     `json:"dist-tags"`
}

func ManifestFromPackageJSON(m manifest) (*entities.PackageVersion, []fields.Email, error) {

	// convert required fields

	name, err := fields.RequiredStringFromString(m.Name)
	if err != nil {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "name",
			Reason: err.Error(),
		}
	}

	var ver string

	for v := range m.Versions {
		ver = v
		break
	}

	if ver == "" {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "version",
			Reason: "no version found",
		}
	}

	version, err := fields.RequiredStringFromString(ver)
	if err != nil {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "version",
			Reason: err.Error(),
		}
	}

	integrity, err := fields.RequiredStringFromString(m.Versions[ver].Dist.Integrity)
	if err != nil {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "integrity",
			Reason: err.Error(),
		}
	}

	shasum, err := fields.RequiredStringFromString(m.Versions[ver].Dist.SHASUM)
	if err != nil {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "shasum",
			Reason: err.Error(),
		}
	}

	// tarball is everthing after "-/" inside m.Versions[ver].Dist.Tarball
	tarball := m.Versions[ver].Dist.Tarball
	tarball = tarball[strings.Index(tarball, "-/")+2:]

	contentType, err := fields.RequiredStringFromString(m.Attachments[tarball].ContentType)
	if err != nil {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "content_type",
			Reason: err.Error(),
		}
	}

	data, err := fields.RequiredStringFromString(m.Attachments[tarball].Data)
	if err != nil {
		return nil, nil, &PackageAdapterManifestConvertFieldError{
			Field:  "data",
			Reason: err.Error(),
		}
	}

	// convert optional fields

	var description *string
	if m.Versions[ver].Description != "" {
		description = new(string)
		*description = m.Versions[ver].Description
	}

	var keywords []fields.RequiredString
	if len(m.Versions[ver].Keywords) > 0 {
		keywords = make([]fields.RequiredString, len(m.Versions[ver].Keywords))
		for i, k := range m.Versions[ver].Keywords {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "keywords",
					Reason: err.Error(),
				}
			}
			keywords[i] = k
		}
	}

	var homepage *fields.Website
	if m.Versions[ver].Homepage != "" {
		h, err := fields.WebsiteFromString(m.Versions[ver].Homepage)
		if err != nil {
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "homepage",
				Reason: err.Error(),
			}
		}
		homepage = &h
	}

	var bugs *fields.Bugs
	if m.Versions[ver].Bugs != nil && (m.Versions[ver].Bugs.URL != "" || m.Versions[ver].Bugs.Email != "") {
		bugsBuilder := fields.BugsBuilder()
		if m.Versions[ver].Bugs.URL != "" {
			bugsBuilder.URL(m.Versions[ver].Bugs.URL)
		}
		if m.Versions[ver].Bugs.Email != "" {
			bugsBuilder.Email(m.Versions[ver].Bugs.Email)
		}
		b, err := bugsBuilder.Build()
		if err != nil {
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "bugs",
				Reason: err.Error(),
			}
		}
		bugs = &b
	}

	var license *string
	if m.Versions[ver].License != "" {
		license = new(string)
		*license = m.Versions[ver].License
	}

	contributersToCheck := make([]fields.Email, len(m.Versions[ver].Contributors))
	contributers := make(fields.MixedAuthors, len(m.Versions[ver].Contributors))
	if len(m.Versions[ver].Contributors) > 0 {
		for i, c := range m.Versions[ver].Contributors {
			if c.Email != nil {
				email, err := fields.EmailFromString(*c.Email)
				if err != nil {
					return nil, nil, &PackageAdapterManifestConvertFieldError{
						Field:  "contributors",
						Reason: err.Error(),
					}
				}
				contributersToCheck[i] = email
			}
		}
	}

	var funding []fields.UrlType

	if m.Versions[ver].Funding != nil {
		funding = make([]fields.UrlType, 0)

		switch f := m.Versions[ver].Funding.(type) {
		case map[string]any:
			fundingBuilder := fields.UrlTypeBuilder()
			if f["url"] != nil {
				fundingBuilder.URL(f["url"].(string))
			}
			if f["type"] != nil {
				fundingBuilder.Type(f["type"].(string))
			}
			fnd, err := fundingBuilder.Build()
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "funding",
					Reason: err.Error(),
				}
			}
			funding = append(funding, fnd)
		case []any:
			for _, v := range f {
				switch v := v.(type) {
				case string:
					fundingBuilder := fields.UrlTypeBuilder()
					fundingBuilder.URL(v)
					fnd, err := fundingBuilder.Build()
					if err != nil {
						return nil, nil, &PackageAdapterManifestConvertFieldError{
							Field:  "funding",
							Reason: err.Error(),
						}
					}
					funding = append(funding, fnd)

				case map[string]any:
					fundingBuilder := fields.UrlTypeBuilder()
					if v["url"] != nil {
						fundingBuilder.URL(v["url"].(string))
					}
					if v["type"] != nil {
						fundingBuilder.Type(v["type"].(string))
					}
					fnd, err := fundingBuilder.Build()
					if err != nil {
						return nil, nil, &PackageAdapterManifestConvertFieldError{
							Field:  "funding",
							Reason: err.Error(),
						}
					}
					funding = append(funding, fnd)
				default:
					return nil, nil, &PackageAdapterManifestConvertFieldError{
						Field:  "funding",
						Reason: fmt.Sprintf("unknown funding entry type %T", v),
					}
				}
			}
		default:
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "funding",
				Reason: fmt.Sprintf("unknown funding type %T", f),
			}
		}
	}

	var files []fields.RequiredString
	if len(m.Versions[ver].Files) > 0 {
		files = make([]fields.RequiredString, len(m.Versions[ver].Files))
		for i, f := range m.Versions[ver].Files {
			f, err := fields.RequiredStringFromString(f)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "files",
					Reason: err.Error(),
				}
			}
			files[i] = f
		}
	}

	var main *fields.RequiredString
	if m.Versions[ver].Main != "" {
		m, err := fields.RequiredStringFromString(m.Versions[ver].Main)
		if err != nil {
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "main",
				Reason: err.Error(),
			}
		}
		main = &m
	}

	var browser *fields.RequiredString
	if m.Versions[ver].Browser != "" {
		b, err := fields.RequiredStringFromString(m.Versions[ver].Browser)
		if err != nil {
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "browser",
				Reason: err.Error(),
			}
		}
		browser = &b
	}

	var bin map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].Bin) > 0 {
		bin = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].Bin {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "bin.key",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "bin.value",
					Reason: err.Error(),
				}
			}
			bin[k] = v
		}
	}

	var man []fields.RequiredString
	if len(m.Versions[ver].Man) > 0 {
		man = make([]fields.RequiredString, len(m.Versions[ver].Man))
		for i, v := range m.Versions[ver].Man {
			m, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "man",
					Reason: err.Error(),
				}
			}
			man[i] = m
		}
	}

	var directories *fields.Directories
	if m.Versions[ver].Directories != nil {
		directoriesBuilder := fields.DirectoriesBuilder()
		if m.Versions[ver].Directories.Bin != "" {
			directoriesBuilder.Bin(m.Versions[ver].Directories.Bin)
		}

		if m.Versions[ver].Directories.Man != "" {
			directoriesBuilder.Man(m.Versions[ver].Directories.Man)
		}

		d, err := directoriesBuilder.Build()
		if err != nil {
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "directories",
				Reason: err.Error(),
			}
		}

		directories = &d
	}

	var repository *fields.Repository
	if m.Versions[ver].Repository != nil && m.Versions[ver].Repository.URL != "" {
		repositoryBuilder := fields.RepositoryBuilder()
		if m.Versions[ver].Repository.Type != nil {
			repositoryBuilder.Type(*m.Versions[ver].Repository.Type)
		}
		if m.Versions[ver].Repository.Directory != nil {
			repositoryBuilder.Directory(*m.Versions[ver].Repository.Directory)
		}
		if m.Versions[ver].Repository.URL != "" {
			repositoryBuilder.URL(m.Versions[ver].Repository.URL)
		}
		r, err := repositoryBuilder.Build()
		if err != nil {
			return nil, nil, &PackageAdapterManifestConvertFieldError{
				Field:  "repository",
				Reason: err.Error(),
			}
		}
		repository = &r
	}

	var scripts map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].Scripts) > 0 {
		scripts = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].Scripts {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "scripts",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "scripts",
					Reason: err.Error(),
				}
			}
			scripts[k] = v
		}
	}

	var config map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].Config) > 0 {
		config = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].Config {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "config.key",
					Reason: err.Error(),
				}
			}

			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "config.value",
					Reason: err.Error(),
				}
			}

			config[k] = v
		}
	}

	var dependencies map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].Dependencies) > 0 {
		dependencies = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].Dependencies {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "dependencies.key",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "dependencies.value",
					Reason: err.Error(),
				}
			}
			dependencies[k] = v
		}
	}

	var devDependencies map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].DevDependencies) > 0 {
		devDependencies = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].DevDependencies {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "devDependencies.key",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "devDependencies.value",
					Reason: err.Error(),
				}
			}
			devDependencies[k] = v
		}
	}

	var peerDependencies map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].PeerDependencies) > 0 {
		peerDependencies = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].PeerDependencies {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "peerDependencies.key",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "peerDependencies.value",
					Reason: err.Error(),
				}
			}
			peerDependencies[k] = v
		}
	}

	var peerDependenciesMeta map[fields.RequiredString]map[fields.RequiredString]any
	if len(m.Versions[ver].PeerDependenciesMeta) > 0 {
		peerDependenciesMeta = make(map[fields.RequiredString]map[fields.RequiredString]any)
		for k, v := range m.Versions[ver].PeerDependenciesMeta {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "peerDependenciesMeta.key",
					Reason: err.Error(),
				}
			}
			peerDependenciesMeta[k] = make(map[fields.RequiredString]any)
			for k2, v2 := range v {
				k2, err := fields.RequiredStringFromString(k2)
				if err != nil {
					return nil, nil, &PackageAdapterManifestConvertFieldError{
						Field:  "peerDependenciesMeta.key2",
						Reason: err.Error(),
					}
				}
				peerDependenciesMeta[k][k2] = v2
			}
		}
	}

	var bundledDependencies []fields.RequiredString
	if len(m.Versions[ver].BundledDependencies) > 0 {
		bundledDependencies = make([]fields.RequiredString, len(m.Versions[ver].BundledDependencies))
		for i, v := range m.Versions[ver].BundledDependencies {
			b, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "bundledDependencies",
					Reason: err.Error(),
				}
			}
			bundledDependencies[i] = b
		}
	}

	var optionalDependencies map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].OptionalDependencies) > 0 {
		optionalDependencies = make(map[fields.RequiredString]fields.RequiredString)
		for k, v := range m.Versions[ver].OptionalDependencies {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "optionalDependencies.key",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "optionalDependencies.value",
					Reason: err.Error(),
				}
			}
			optionalDependencies[k] = v
		}
	}

	var engines map[fields.RequiredString]fields.RequiredString
	if len(m.Versions[ver].Engines) > 0 {
		engines = make(map[fields.RequiredString]fields.RequiredString)

		for k, v := range m.Versions[ver].Engines {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "engines.key",
					Reason: err.Error(),
				}
			}
			v, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "engines.value",
					Reason: err.Error(),
				}
			}
			engines[k] = v
		}
	}

	var os []fields.RequiredString
	if len(m.Versions[ver].OS) > 0 {
		os = make([]fields.RequiredString, len(m.Versions[ver].OS))
		for i, v := range m.Versions[ver].OS {
			o, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "os",
					Reason: err.Error(),
				}
			}
			os[i] = o
		}
	}

	var cpu []fields.RequiredString
	if len(m.Versions[ver].CPU) > 0 {
		cpu = make([]fields.RequiredString, len(m.Versions[ver].CPU))
		for i, v := range m.Versions[ver].CPU {
			c, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "cpu",
					Reason: err.Error(),
				}
			}
			cpu[i] = c
		}
	}

	private := false
	if m.Versions[ver].Private != nil {
		private = *m.Versions[ver].Private
	}

	var publishConfig map[fields.RequiredString]any
	if len(m.Versions[ver].PublishConfig) > 0 {
		publishConfig = make(map[fields.RequiredString]any)
		for k, v := range m.Versions[ver].PublishConfig {
			k, err := fields.RequiredStringFromString(k)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "publishConfig.key",
					Reason: err.Error(),
				}
			}
			publishConfig[k] = v
		}
	}

	var workspaces []fields.RequiredString
	if len(m.Versions[ver].Workspaces) > 0 {
		workspaces = make([]fields.RequiredString, len(m.Versions[ver].Workspaces))
		for i, v := range m.Versions[ver].Workspaces {
			w, err := fields.RequiredStringFromString(v)
			if err != nil {
				return nil, nil, &PackageAdapterManifestConvertFieldError{
					Field:  "workspaces",
					Reason: err.Error(),
				}
			}
			workspaces[i] = w
		}
	}

	var length int
	if m.Attachments[tarball].Length != 0 {
		length = m.Attachments[tarball].Length
	}

	var readme *string

	if m.Versions[ver].Readme != "" {
		readme = new(string)
		*readme = m.Versions[ver].Readme
	}

	return &entities.PackageVersion{
		Name:            name,
		Version:         version,
		Description:     description,
		Homepage:        homepage,
		Repository:      repository,
		Dependencies:    dependencies,
		DevDependencies: devDependencies,
		Scripts:         scripts,
		License:         license,
		Integrity:       integrity,
		SHASUM:          shasum,
		ContentType:     contentType,
		Data:            data,
		Length:          length,
		Readme:          readme,
		Contributors:    contributers,

		Keywords:             keywords,
		Bugs:                 bugs,
		Funding:              funding,
		Files:                files,
		Main:                 main,
		Browser:              browser,
		Bin:                  bin,
		Man:                  man,
		Directories:          directories,
		PeerDependencies:     peerDependencies,
		PeerDependenciesMeta: peerDependenciesMeta,
		BundledDependencies:  bundledDependencies,
		OptionalDependencies: optionalDependencies,
		Engines:              engines,
		OS:                   os,
		CPU:                  cpu,
		Private:              &private,
		PublishConfig:        publishConfig,
		Workspaces:           workspaces,
	}, contributersToCheck, nil
}

// errors

type PackageAdapterManifestConvertFieldError struct {
	Field  string
	Reason string
}

func (e *PackageAdapterManifestConvertFieldError) Error() string {
	return fmt.Sprintf("error converting package manifest field %s: %s", e.Field, e.Reason)
}
