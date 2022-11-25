package main

var servers = []Server{
	{URL: "https://norrebro.space"},
	{URL: "https://turingfesten.dk"},
	{URL: "https://mstdn.dk"},
	{URL: "https://helvede.net"},
	{URL: "https://expressional.social"},
	{URL: "https://uddannelse.social"},
	{URL: "https://social.data.coop", ForceCategory: PrivateCategory},
}

var (
	OpenCategory = &Category{
		Name:        "Åbne servere",
		Description: "Servere du kan blive medlem af med det samme, uden godkendelse.",
		Servers:     make([]ServerResponse, 0),
	}

	OpenCategoryWithApproval = &Category{
		Name:        "Server med manuel godkendelse",
		Description: "Servere der ønsker at godkende dit medlemsskab manuelt efter oprettelse.",
		Servers:     make([]ServerResponse, 0),
	}

	ClosedCategory = &Category{
		Name:        "Accepterer ikke nye medlemmer",
		Description: "Servere der ikke længere accepterer nye medlemmer",
		Servers:     make([]ServerResponse, 0),
	}

	PrivateCategory = &Category{
		Name:        "Forening / private",
		Description: "Servere der kræver medlemsskab af en forening eller andet specielt for at få en bruger.",
		Servers:     make([]ServerResponse, 0),
	}
)
