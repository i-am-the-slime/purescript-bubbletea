{ name = "bubbletea"
, dependencies = 
  [ "effect"
  , "prelude"
  , "maybe"
  , "arrays"
  , "functions"
  ]
, packages = https://github.com/purescript/package-sets/releases/download/psc-0.15.15-20240701/packages.dhall
    sha256:c0ad6f3f08ab980e9cacfd6e1fc96f111fd48aad70b9f951dde7bab5e2dc954e
, sources = [ "src/**/*.purs" ]
, backend = "psgo"
}