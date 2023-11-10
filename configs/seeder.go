package configs

import (
	"errors"
	"os"
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"time"

	"gorm.io/gorm"
)

func checkEnumIsExist(db *gorm.DB, key string) bool {
	var count int64
	if err := db.Table("pg_type").Where("typname = ?", key).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func SeedEnum(db *gorm.DB) error {
	if !checkEnumIsExist(db, "allergens_enum") {
		return db.Exec(
			`CREATE TYPE allergens_enum 
			AS ENUM (
				'` + string(consttypes.A_FOOD) + `',
				'` + string(consttypes.A_MEDICAL) + `',
				'` + string(consttypes.A_ENVIRONMENTAL) + `',
				'` + string(consttypes.A_CONTACT) + `'
			);`,
		).Error
	}

	return nil
}

func SeedAdminCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Admin{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			time, err := time.Parse(consttypes.DATEFORMAT, "2000-10-20")
			if err != nil {
				return err
			}

			admins := []*models.Admin{
				{
					User: models.User{
						Email:    os.Getenv("ADMIN_EMAIL"),
						Password: os.Getenv("ADMIN_PASSWORD"),
						Role:     consttypes.UR_ADMIN,
					},
					FirstName:   os.Getenv("ADMIN_FIRSTNAME"),
					LastName:    os.Getenv("ADMIN_LASTNAME"),
					Gender:      consttypes.G_MALE,
					DateOfBirth: time,
				},
			}

			db.Create(admins)
		}
	}

	return nil
}

func SeedAllergyData(db *gorm.DB) error {
	// * source: https://en.wikipedia.org/wiki/List_of_allergens
	// TODO - needs to continue seed the other allergen type
	if db.Migrator().HasTable(&models.Allergy{}) {
		if err := db.First(&models.Allergy{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			allergies := []*models.Allergy{
				// ! start of food allergen
				{
					Name:        "Milk",
					Description: "A milk allergy, also known as a dairy allergy, is an adverse immune system response to one or more proteins found in cow's milk. It is different from lactose intolerance, which is a non-immune digestive disorder where the body has difficulty digesting lactose, a sugar found in milk. A milk allergy is an immune system disorder and can be more severe.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Balsam of Peru",
					Description: "Balsam of Peru (Myroxylon pereirae) is a natural resin derived from certain trees native to Central America and South America. It is commonly used as a fragrance and flavoring agent in various products, including perfumes, cosmetics, toiletries, and food items. However, some individuals can develop an allergy to Balsam of Peru, which can lead to allergic contact dermatitis and other symptoms.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Buckwheat",
					Description: "A buckwheat allergy is an adverse immune response to proteins found in buckwheat, a grain-like seed often used in various foods and dishes. It is a relatively uncommon food allergy but can cause a range of allergic reactions in individuals who are sensitized to buckwheat proteins.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Celery",
					Description: "A celery allergy is an adverse immune response to proteins found in celery, a common vegetable used in various culinary dishes and food products. It is considered one of the more prevalent allergies to vegetables and can cause a range of allergic reactions in individuals who are sensitized to celery proteins.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Egg",
					Description: "An egg allergy is an adverse immune response to proteins found in eggs, most commonly the proteins in egg whites but sometimes also in egg yolks. It is one of the most common food allergies, particularly in children, but some individuals may outgrow it as they get older.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Fish",
					Description: "A fish allergy is an adverse immune response to proteins found in fish, typically marine or saltwater fish like salmon, tuna, and cod. It is one of the most common food allergies and can cause a range of allergic reactions in individuals who are sensitized to fish proteins.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Fruit",
					Description: "A fruit allergy is an adverse immune response to proteins found in various fruits. It is a relatively common food allergy and can cause a range of allergic reactions in individuals who are sensitized to specific fruit proteins.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Garlic",
					Description: "Garlic allergy or allergic contact dermatitis to garlic is a common inflammatory skin condition caused by contact with garlic oil or dust. It mostly affects people who cut and handle fresh garlic, such as chefs, and presents on the tips of the thumb, index and middle fingers of the non-dominant hand. The affected fingertips show an asymmetrical pattern of fissure as well as thickening and shedding of the outer skin layers, which may progress to second- or third-degree burn of injured skin.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Oats",
					Description: "Oat sensitivity represents a sensitivity to the proteins found in oats, Avena sativa. Sensitivity to oats can manifest as a result of allergy to oat seed storage proteins either inhaled or ingested.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Maize",
					Description: "Corn allergy is a very rare food allergy. People with a true IgE-mediated allergy to corn develop symptoms such as swelling or hives when they eat corn or foods that contain corn. The allergy can be difficult to manage due to many food and non-food products that contain various forms of corn, such as corn starch and modified food starch.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Mustard",
					Description: "A mustard allergy is an adverse immune response to proteins found in mustard seeds. It is a relatively uncommon food allergy but can cause a range of allergic reactions in individuals who are sensitized to mustard proteins.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Peanut",
					Description: "Peanut allergy is a type of food allergy to peanuts. It is different from tree nut allergies, because peanuts are legumes and not true nuts. Physical symptoms of allergic reaction can include itchiness, hives, swelling, eczema, sneezing, asthma attack, abdominal pain, drop in blood pressure, diarrhea, and cardiac arrest. Anaphylaxis may occur. Those with a history of asthma are more likely to be severely affected.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Poultry Meat",
					Description: "Poultry meat allergy is a rare food allergy in humans caused by consumption of poultry meat whereby the body triggers an immune reaction and becomes overloaded with immunoglobulin E (IgE) antibodies.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Red Meat",
					Description: "A red meat allergy, also known as alpha-gal allergy, is an adverse immune response to a carbohydrate molecule called galactose-alpha-1,3-galactose (alpha-gal), which is found in red meat, particularly from mammals like cows, pigs, and sheep. It is considered a relatively rare but intriguing type of food allergy with distinctive characteristics.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Rice",
					Description: "Rice allergy is a type of food allergy. People allergic to rice react to various rice proteins after they eat rice or breathe the steam from cooking rice. Although some reactions might lead to severe health problems, doctors can diagnose rice allergy with many methods and help allergic people to avoid reactions.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Sesame",
					Description: "A sesame allergy is an adverse immune response to proteins found in sesame seeds, a common ingredient in various foods and cuisines. It is considered one of the more prevalent food allergies and can cause a range of allergic reactions in individuals who are sensitized to sesame proteins.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Shellfish",
					Description: "Shellfish allergy is among the most common food allergies. 'Shellfish' is a colloquial and fisheries term for aquatic invertebrates used as food, including various species of molluscs such as clams, mussels, oysters and scallops, crustaceans such as shrimp, lobsters and crabs, and cephalopods such as squid and octopus.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Soy",
					Description: "Soy allergy is a type of food allergy. It is a hypersensitivity to ingesting compounds in soy, causing an overreaction of the immune system, typically with physical symptoms, such as gastrointestinal discomfort, respiratory distress, or a skin reaction. Soy is among the eight most common foods inducing allergic reactions in children and adults. It has a prevalence of about 0.3% in the general population.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Sulfites",
					Description: "Sulfite sensitivity, also known as sulfite allergy or sulfite intolerance, refers to an adverse reaction to sulfites, which are sulfur-based compounds commonly used as food preservatives, primarily in dried fruits, wine, and some processed foods. It is important to note that sulfite sensitivity is different from a true allergy in the immunological sense, as it does not involve the immune system's IgE-mediated response, but it can still lead to unpleasant symptoms.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Tartrazine",
					Description: "Tartrazine, also known as Yellow 5, is a synthetic yellow food coloring and dye used in various food and beverage products. While not strictly an allergy, some individuals may be sensitive to tartrazine, experiencing adverse reactions upon consumption.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Tree nut",
					Description: "A tree nut allergy is a hypersensitivity to dietary substances from tree nuts and edible tree seeds causing an overreaction of the immune system which may lead to severe physical symptoms. Tree nuts include almonds, Brazil nuts, cashews, chestnuts, filberts/hazelnuts, macadamia nuts, pecans, pistachios, shea nuts and walnuts.",
					Allergens:   consttypes.A_FOOD,
				},
				{
					Name:        "Wheat",
					Description: "Wheat allergy is an allergy to wheat which typically presents itself as a food allergy, but can also be a contact allergy resulting from occupational exposure. Like all allergies, wheat allergy involves immunoglobulin E and mast cell response. Typically the allergy is limited to the seed storage proteins of wheat.",
					Allergens:   consttypes.A_FOOD,
				},
				// ! start of medical allergen
				{
					Name:        "Balsam of Peru",
					Description: "Balsam of Peru is a resinous substance that comes from the Myroxylon balsamum tree, which is native to Central America. Despite its name, it is not a true balsam. The resin is extracted from the tree trunk and branches. Balsam of Peru is often used in perfumes, flavorings, and various cosmetic and medicinal products.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Tetracycline",
					Description: "Tetracycline, sold under various brand names, is an oral antibiotic in the tetracyclines family of medications, used to treat a number of infections, including acne, cholera, brucellosis, plague, malaria, and syphilis.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Dilantin",
					Description: "Phenytoin (PHT), sold under the brand name Dilantin among others, is an anti-seizure medication. It is useful for the prevention of tonic-clonic seizures and focal seizures, but not absence seizures. The intravenous form, fosphenytoin, is used for status epilepticus that does not improve with benzodiazepines.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Tegretol (carbamazepine)",
					Description: "Carbamazepine, sold under the brand name Tegretol among others, is an anticonvulsant medication used in the treatment of epilepsy and neuropathic pain. It is used as an adjunctive treatment in schizophrenia along with other medications and as a second-line agent in bipolar disorder. Carbamazepine appears to work as well as phenytoin and valproate for focal and generalized seizures.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Penicillin",
					Description: "Carbamazepine, sold under the brand name Tegretol among others, is an anticonvulsant medication used in the treatment of epilepsy and neuropathic pain. It is used as an adjunctive treatment in schizophrenia along with other medications and as a second-line agent in bipolar disorder. Carbamazepine appears to work as well as phenytoin and valproate for focal and generalized seizures.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Cephalosporins",
					Description: "The cephalosporins are a class of β-lactam antibiotics originally derived from the fungus Acremonium, which was previously known as Cephalosporium.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Sulfonamides",
					Description: "Sulfonamide is a functional group that is the basis of several groups of drugs, which are called sulphonamides, sulfa drugs or sulpha drugs. The original antibacterial sulfonamides are synthetic (nonantibiotic) antimicrobial agents that contain the sulfonamide group. Some sulfonamides are also devoid of antibacterial activity, e.g., the anticonvulsant sultiame.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Non-steroidal anti-inflammatories",
					Description: "Non-steroidal anti-inflammatory drugs (NSAID) are members of a therapeutic drug class which reduces pain, decreases inflammation, decreases fever, and prevents blood clots. Side effects depend on the specific drug, its dose and duration of use, but largely include an increased risk of gastrointestinal ulcers and bleeds, heart attack, and kidney disease.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Intravenous contrast dye",
					Description: "Radiocontrast agents are substances used to enhance the visibility of internal structures in X-ray-based imaging techniques such as computed tomography, projectional radiography, and fluoroscopy. Radiocontrast agents are typically iodine, or more rarely barium sulfate. The contrast agents absorb external X-rays, resulting in decreased exposure on the X-ray detector.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Local anesthetics",
					Description: "A local anesthetic (LA) is a medication that causes absence of all sensation in a specific body part without loss of consciousness, as opposed to a general anesthetic, which eliminates all sensation in the entire body and causes unconsciousness.",
					Allergens:   consttypes.A_MEDICAL,
				},
				// ! start of contact allergen
				{
					Name:        "Dimethylaminopropylamine",
					Description: "Dimethylaminopropylamine (DMAPA) is a diamine used in the preparation of some surfactants, such as cocamidopropyl betaine which is an ingredient in many personal care products including soaps, shampoos, and cosmetics. BASF, a major producer, claims that DMAPA-derivatives do not sting the eyes and makes a fine-bubble foam, making it appropriate in shampoos.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Latex",
					Description: "Latex is an emulsion of polymer microparticles in water. Latexes are found in nature, but synthetic latexes are common as well.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Paraphenylenediamine",
					Description: "p-Phenylenediamine (PPD) is an organic compound with the formula C6H4(NH2)2. This derivative of aniline is a white solid, but samples can darken due to air oxidation. It is mainly used as a component of engineering polymers and composites like kevlar. It is also an ingredient in hair dyes and is occasionally used as a substitute for henna.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Glyceryl monothioglycolate",
					Description: "A permanent wave, commonly called a perm or permanent, is a hairstyle consisting of waves or curls set into the hair. The curls may last a number of months, hence the name.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Toluenesulfonamide formaldehyde ",
					Description: "There are three isomers of toluidine, which are organic compounds. These isomers are o-toluidine, m-toluidine, and p-toluidine, with the prefixed letter abbreviating, respectively, ortho; meta; and para. All three are aryl amines whose chemical structures are similar to aniline except that a methyl group is substituted onto the benzene ring. The difference between these three isomers is the position where the methyl group (-CH3) is bonded to the ring relative to the amino functional group (-NH2); see illustration of the chemical structures below.",
					Allergens:   consttypes.A_MEDICAL,
				},
			}

			db.Create(allergies)
		}
	}

	return nil
}
