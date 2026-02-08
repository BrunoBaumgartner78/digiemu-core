package domain

// Meaning represents optional structured context attached to a Unit/Version.
// Fields are intentionally optional (omitempty) to keep compatibility.
type Meaning struct {
	SchemaVersion string             `json:"schema_version"`
	Title         string             `json:"title,omitempty"`
	Purpose       string             `json:"purpose,omitempty"`
	Scope         *MeaningScope      `json:"scope,omitempty"`
	Claims        []MeaningClaim     `json:"claims,omitempty"`
	Sources       []MeaningSource    `json:"sources,omitempty"`
	Provenance    *MeaningProvenance `json:"provenance,omitempty"`
	Integrity     *MeaningIntegrity  `json:"integrity,omitempty"`
}

type MeaningScope struct {
	Audience     []string          `json:"audience,omitempty"`
	Jurisdiction []string          `json:"jurisdiction,omitempty"`
	Locale       []string          `json:"locale,omitempty"`
	Timeframe    *MeaningTimeframe `json:"timeframe,omitempty"`
}

type MeaningTimeframe struct {
	ValidFrom  string `json:"valid_from,omitempty"`
	ValidUntil string `json:"valid_until,omitempty"`
}

type MeaningClaim struct {
	Text     string   `json:"text,omitempty"`
	Strength string   `json:"strength,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type MeaningSource struct {
	ID    string              `json:"id,omitempty"`
	Type  string              `json:"type,omitempty"`
	Ref   string              `json:"ref,omitempty"`
	Quote *MeaningSourceQuote `json:"quote,omitempty"`
}

type MeaningSourceQuote struct {
	Snippet string `json:"snippet,omitempty"`
	Locator string `json:"locator,omitempty"`
}

type MeaningProvenance struct {
	Author    string `json:"author,omitempty"`
	Org       string `json:"org,omitempty"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type MeaningIntegrity struct {
	NarrativeID   string   `json:"narrative_id,omitempty"`
	Supersedes    string   `json:"supersedes,omitempty"`
	ConflictsWith []string `json:"conflicts_with,omitempty"`
}

// MeaningRef is used in export manifests to reference a meaning document.
type MeaningRef struct {
	MeaningHash   string `json:"meaning_hash"`
	MeaningPath   string `json:"meaning_path,omitempty"`
	SchemaVersion string `json:"schema_version,omitempty"`
}

// MeaningHash is the hex-encoded sha256 of the canonicalized meaning.json
type MeaningHash string
