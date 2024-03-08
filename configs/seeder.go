package configs

import (
	"errors"
	"fmt"
	"os"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"

	"project-skbackend/packages/utils/utlogger"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	uuidval, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
)

func getGlobalHashedPassword(password string) string {
	hash, err := helper.HashPassword(password)
	if err != nil {
		utlogger.LogError(err)
		return ""
	}

	return hash
}

func checkEnumIsExist(db *gorm.DB, key string) bool {
	var count int64
	if err := db.Table("pg_type").Where("typname = ?", key).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func createEnum(db *gorm.DB, enumname string, enumvalues ...string) error {
	if !checkEnumIsExist(db, enumname) {
		values := make([]string, len(enumvalues))
		for i, v := range enumvalues {
			values[i] = string(v)
		}

		query := fmt.Sprintf("CREATE TYPE %s AS ENUM ('%s');", enumname, strings.Join(values, "','"))
		err := db.Exec(query).Error
		if err != nil {
			utlogger.LogError(err)
			return err
		}
	}
	return nil
}

func SeedAllergensEnum(db *gorm.DB) error {
	return createEnum(db,
		"allergens_enum",
		consttypes.A_FOOD.String(),
		consttypes.A_MEDICAL.String(),
		consttypes.A_ENVIRONMENTAL.String(),
		consttypes.A_CONTACT.String(),
	)
}

func SeedUserRoleEnum(db *gorm.DB) error {
	return createEnum(db,
		"user_role_enum",
		consttypes.UR_ADMIN.String(),
		consttypes.UR_CAREGIVER.String(),
		consttypes.UR_MEMBER.String(),
		consttypes.UR_ORGANIZATION.String(),
		consttypes.UR_PARTNER.String(),
		consttypes.UR_PATRON.String(),
		consttypes.UR_USER.String(),
	)
}

func SeedGenderEnum(db *gorm.DB) error {
	return createEnum(db,
		"gender_enum",
		consttypes.G_MALE.String(),
		consttypes.G_FEMALE.String(),
		consttypes.G_OTHER.String(),
	)
}

func SeedMealStatusEnum(db *gorm.DB) error {
	return createEnum(db,
		"meal_status_enum",
		consttypes.MS_ACTIVE.String(),
		consttypes.MS_INACTIVE.String(),
		consttypes.MS_OUTOFSTOCK.String(),
	)
}

func SeedDonationStatusEnum(db *gorm.DB) error {
	return createEnum(db,
		"donation_status_enum",
		consttypes.DS_ACCEPTED.String(),
		consttypes.DS_REJECTED.String(),
	)
}

func SeedImageTypeEnum(db *gorm.DB) error {
	return createEnum(db,
		"image_type_enum",
		consttypes.IT_PROFILE.String(),
		consttypes.IT_MEAL.String(),
	)
}

func SeedPatronTypeEnum(db *gorm.DB) error {
	return createEnum(db,
		"patron_type_enum",
		consttypes.PT_ORGANIZATION.String(),
		consttypes.PT_PERSONAL.String(),
	)
}

func SeedOrganizationTypeEnum(db *gorm.DB) error {
	return createEnum(db,
		"organization_type_enum",
		consttypes.OT_NURSINGHOME.String(),
	)
}

func SeedAdminCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Admin{}) {
		if err := db.First(&models.Admin{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			admins := []*models.Admin{
				{
					User: models.User{
						Email:    os.Getenv("ADMIN_EMAIL"),
						Password: getGlobalHashedPassword(os.Getenv("ADMIN_PASSWORD")),
						Role:     consttypes.UR_ADMIN,
					},
					FirstName:   os.Getenv("ADMIN_FIRSTNAME"),
					LastName:    os.Getenv("ADMIN_LASTNAME"),
					Gender:      consttypes.G_MALE,
					DateOfBirth: customs.CDT_DATE{Time: consttypes.DateNow},
				},
			}

			db.Create(admins)
		}
	}

	return nil
}

func SeedCaregiverCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Caregiver{}) {
		if err := db.First(&models.Caregiver{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			caregivers := []*models.Caregiver{
				{
					User: models.User{
						Email:    "caregiver@test.com",
						Password: getGlobalHashedPassword("password"),
						Role:     consttypes.UR_CAREGIVER,
					},
					FirstName:   "Care",
					LastName:    "Giver",
					Gender:      consttypes.G_FEMALE,
					DateOfBirth: customs.CDT_DATE{Time: consttypes.DateNow},
				},
			}

			db.Create(caregivers)
		}
	}

	return nil
}

func SeedOrganizationCredential(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Organization{}) {
		if err := db.First(&models.Organization{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			organizations := []*models.Organization{
				{
					User: models.User{
						Email:    "organization@test.com",
						Password: getGlobalHashedPassword("password"),
						Role:     consttypes.UR_ORGANIZATION,
					},
					Type: consttypes.OT_NURSINGHOME,
					Name: "Nursing Home",
				},
			}

			db.Create(organizations)
		}
	}

	return nil
}

func SeedPartnerCredential(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Partner{}) {
		if err := db.First(&models.Partner{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			partners := []*models.Partner{
				{
					Model: helper.Model{
						ID: uuidval,
					},
					User: models.User{
						Email:    "partner@test.com",
						Password: getGlobalHashedPassword("password"),
						Role:     consttypes.UR_PARTNER,
					},
					Name: "Partner",
				},
			}

			err = db.Create(partners).Error
			if err != nil {
				utlogger.LogError(err)
				return err
			}
		}
	}

	return nil
}

func SeedMealData(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.Meal{}) && db.Migrator().HasTable(&models.Partner{}) {
		if err := db.First(&models.Meal{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var illness models.Illness
			db.First(&illness, uuidval)

			var allergy models.Allergy
			db.First(&allergy, uuidval)

			var partner models.Partner
			db.First(&partner, uuidval)

			meals := []*models.Meal{
				{
					Illnesses: []*models.MealIllness{
						{
							Illness: illness,
						},
					},
					Allergies: []*models.MealAllergy{
						{
							Allergy: allergy,
						},
					},
					Partner:     partner,
					Name:        "Nasi Goyeng",
					Status:      consttypes.MS_ACTIVE,
					Description: "Nasi goyeng adalah masakan Indonesia yang populer, terkenal karena rasa yang kaya dan beragam. Ini adalah hidangan nasi yang digoreng dengan bumbu-bumbu khas Indonesia dan seringkali ditambahkan dengan berbagai jenis bahan tambahan seperti daging, ayam, udang, telur, sayuran, dan rempah-rempah.",
				},
			}

			err = db.Create(meals).Error
			utlogger.LogError(err)

		}
	}

	return nil
}

func SeedAllergyData(db *gorm.DB) error {
	// * source: https://en.wikipedia.org/wiki/List_of_allergens
	if db.Migrator().HasTable(&models.Allergy{}) {
		if err := db.First(&models.Allergy{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			allergies := []*models.Allergy{
				// ! start of food allergen
				{
					Model: helper.Model{
						ID: uuidval,
					},
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
					Description: "Glyceryl monothioglycolate is a chemical compound primarily used in hair products, particularly in hair waving or straightening solutions. It's an ingredient commonly found in hair relaxers or perms.",
					Allergens:   consttypes.A_MEDICAL,
				},
				{
					Name:        "Toluenesulfonamide formaldehyde ",
					Description: "There are three isomers of toluidine, which are organic compounds. These isomers are o-toluidine, m-toluidine, and p-toluidine, with the prefixed letter abbreviating, respectively, ortho; meta; and para. All three are aryl amines whose chemical structures are similar to aniline except that a methyl group is substituted onto the benzene ring. The difference between these three isomers is the position where the methyl group (-CH3) is bonded to the ring relative to the amino functional group (-NH2); see illustration of the chemical structures below.",
					Allergens:   consttypes.A_MEDICAL,
				},
				// ! start of environmental allergen
				{
					Name:        "Pollen",
					Description: "Pollens are microspores from trees, grass or weeds that appear as a fine dust. Pollen may be many colors, including yellow, white, red or brown. Plants release pollen to fertilize other plants for reproduction. Pollen levels are usually highest in the morning. Pollen levels increase on warm, windy days.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
				{
					Name:        "Molds",
					Description: "Molds are tiny fungi (singular, fungus). They have spores that float in the air. Mold is common in damp areas with little or no airflow. These areas may include your basement, kitchen or bathroom. Mold also grows outdoors in leaf piles, grass, mulch, hay or under mushrooms. Mold spore levels are highest during hot, humid weather.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
				{
					Name:        "Pet dander and saliva (spit)",
					Description: "Pet dander is tiny scales from your pet's skin, hair or feathers. Your pet's sweat glands secrete proteins through their skin, which collect in their skin and fur and may cause an allergic reaction. Your pet's spit (saliva) also contains these proteins.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
				{
					Name:        "Dust mites",
					Description: "Dust mites are tiny, eight-legged relatives of spiders. They're too small to see with your eyes. They live on bedding, mattresses, carpets, curtains and upholstered (fabric) furniture. They feed on the dead skin cells that you and your pets shed. Dust mites live on every continent except Antarctica, but they thrive in hot, humid environments. They don't bite you. Breathing in the proteins from their urine (pee), feces (poop) and dead bodies may cause allergic reactions.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
				{
					Name:        "Cockroaches",
					Description: "Cockroaches are reddish-brown or black insects that are 1.5 to 2 inches long. Male cockroaches have two pairs of wings. Many female cockroaches don't have wings. If they have wings, they aren't strong enough to allow flight (vestigial wings). The proteins in their poop, spit, eggs and dead body parts may cause allergic reactions.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
				{
					Name:        "Smoke",
					Description: "Smoke of any kind can trigger a non-IgE reaction. The chemicals in these products can cause irritation that's similar to an allergic reaction. Examples include tobacco product smoke — including cigarettes, vapes and cigars — and marijuana and scented candle smoke.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
				{
					Name:        "Dust",
					Description: "Dust is a combination of tiny particles of matter. Dust may include dead skin cells, hair, pollen, clothing fibers, dust mites, dead insect pieces, dirt, bacteria and tiny pieces of plastic.",
					Allergens:   consttypes.A_ENVIRONMENTAL,
				},
			}

			db.Create(allergies)
		}
	}

	return nil
}

func SeedIllnessData(db *gorm.DB) error {
	// * source: https://chat.openai.com
	if db.Migrator().HasTable(&models.Illness{}) {
		if err := db.First(&models.Illness{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			illnesses := []*models.Illness{
				{
					Model: helper.Model{
						ID: uuidval,
					},
					Name:        "Covid",
					Description: "Coronavirus disease (COVID-19) is an infectious disease caused by the SARS-CoV-2 virus.",
				},
				{
					Name:        "Influenza",
					Description: "Influenza, commonly known as the flu, is a contagious respiratory illness caused by influenza viruses that infect the nose, throat, and sometimes the lungs.",
				},
				{
					Name:        "Common Cold",
					Description: "The common cold is a viral infection of the upper respiratory tract, primarily caused by rhinoviruses, leading to symptoms like sneezing, runny or stuffy nose, sore throat, and coughing.",
				},
				{
					Name:        "Pneumonia",
					Description: "Pneumonia is an infection that inflames the air sacs in one or both lungs, often caused by bacteria, viruses, or fungi, leading to symptoms like cough, fever, and difficulty breathing.",
				},
				{
					Name:        "Bronchitis",
					Description: "Bronchitis is the inflammation of the lining of the bronchial tubes, often due to viral infections, causing symptoms such as cough, chest discomfort, and production of mucus.",
				},
				{
					Name:        "Urinary Tract Infection (UTI)",
					Description: "A urinary tract infection is an infection in any part of the urinary system, commonly caused by bacteria and leading to symptoms like frequent urination, pain or burning sensation during urination, and abdominal pain.",
				},
				{
					Name:        "Gastroenteritis",
					Description: "Gastroenteritis is an inflammation of the stomach and intestines, commonly caused by viral or bacterial infection, resulting in symptoms like diarrhea, vomiting, abdominal cramps, and nausea.",
				},
				{
					Name:        "Migraine",
					Description: "A migraine is a type of headache characterized by throbbing pain, often accompanied by nausea, vomiting, and sensitivity to light and sound.",
				},
				{
					Name:        "Asthma",
					Description: "Asthma is a chronic respiratory condition that causes airways to become inflamed and narrowed, leading to difficulty breathing, wheezing, coughing, and chest tightness.",
				},
				{
					Name:        "Diabetes",
					Description: "Diabetes is a chronic condition characterized by high levels of sugar (glucose) in the blood, resulting from inadequate insulin production or ineffective use of insulin by the body.",
				},
				{
					Name:        "Hypertension (High Blood Pressure)",
					Description: "Hypertension, or high blood pressure, is a condition where the force of blood against the artery walls is consistently too high, potentially leading to various health problems such as heart disease, stroke, and kidney issues.",
				},
				{
					Name:        "Arthritis",
					Description: "Arthritis is a broad term referring to joint inflammation, causing pain, stiffness, and swelling. There are many types of arthritis, including osteoarthritis and rheumatoid arthritis.",
				},
				{
					Name:        "Osteoporosis",
					Description: "Osteoporosis is a condition characterized by weakened bones, increasing the risk of fractures. It often occurs due to bone loss, making bones brittle and prone to breakage.",
				},
				{
					Name:        "Depression",
					Description: "Depression is a mental health disorder causing persistent feelings of sadness, loss of interest in activities, changes in appetite or sleep patterns, and difficulty concentrating or making decisions.",
				},
				{
					Name:        "Anxiety Disorders",
					Description: "Anxiety disorders involve excessive worry, fear, or apprehension, leading to symptoms such as restlessness, irritability, difficulty controlling feelings of worry, and physical symptoms like increased heart rate.",
				},
				{
					Name:        "Chronic Obstructive Pulmonary Disease (COPD)",
					Description: "COPD is a group of lung diseases, including chronic bronchitis and emphysema, characterized by breathing difficulties, coughing, wheezing, and tightness in the chest.",
				},
				{
					Name:        "Eczema (Atopic Dermatitis)",
					Description: "Eczema is a skin condition that causes itchy, inflamed patches, often appearing on the face, hands, elbows, or knees. It can be triggered by various factors and tends to be chronic.",
				},
				{
					Name:        "Hyperthyroidism",
					Description: "Hyperthyroidism occurs when the thyroid gland produces too much thyroid hormone, leading to symptoms such as weight loss, rapid heartbeat, sweating, anxiety, and fatigue.",
				},
				{
					Name:        "Hypothyroidism",
					Description: "Hypothyroidism is a condition where the thyroid gland does not produce enough thyroid hormone, causing symptoms like fatigue, weight gain, dry skin, constipation, and depression.",
				},
				{
					Name:        "Gastroesophageal Reflux Disease (GERD)",
					Description: "GERD is a digestive disorder where stomach acid frequently flows back into the esophagus, causing symptoms like heartburn, chest pain, difficulty swallowing, and regurgitation.",
				},
				{
					Name:        "Chronic Kidney Disease (CKD)",
					Description: "Chronic kidney disease involves the gradual loss of kidney function over time. It may lead to complications such as fluid retention, electrolyte imbalances, and anemia.",
				},
				{
					Name:        "Celiac Disease",
					Description: "Celiac disease is an autoimmune disorder where ingestion of gluten leads to damage in the small intestine, resulting in symptoms like diarrhea, abdominal pain, fatigue, and malnutrition.",
				},
				{
					Name:        "Osteoarthritis",
					Description: "Osteoarthritis is a degenerative joint disease characterized by the breakdown of joint cartilage and underlying bone, leading to symptoms like joint pain, stiffness, and reduced mobility.",
				},
				{
					Name:        "Rheumatoid Arthritis",
					Description: "Rheumatoid arthritis is an autoimmune disorder that causes joint inflammation, pain, stiffness, and swelling. It can affect various organs and lead to deformities in severe cases.",
				},
				{
					Name:        "Fibromyalgia",
					Description: "Fibromyalgia is a chronic disorder characterized by widespread musculoskeletal pain, fatigue, sleep disturbances, and tenderness in specific areas of the body.",
				},
				{
					Name:        "Migraine",
					Description: "A migraine is a severe headache that can cause intense throbbing or pulsing sensation, often accompanied by nausea, vomiting, and sensitivity to light and sound.",
				},
				{
					Name:        "Obsessive-Compulsive Disorder (OCD)",
					Description: "OCD is a mental health condition characterized by repetitive, intrusive thoughts (obsessions) and behaviors or rituals (compulsions) done to alleviate anxiety or distress.",
				},
				{
					Name:        "Post-Traumatic Stress Disorder (PTSD)",
					Description: "PTSD is a mental health condition triggered by experiencing or witnessing a traumatic event. Symptoms may include flashbacks, nightmares, severe anxiety, and emotional numbness.",
				},
				{
					Name:        "Chronic Fatigue Syndrome (CFS)",
					Description: "CFS, also known as myalgic encephalomyelitis (ME), is a complex disorder characterized by extreme fatigue that doesn't improve with rest, often accompanied by other symptoms like muscle pain, impaired memory, and sleep issues.",
				},
				{
					Name:        "Gout",
					Description: "Gout is a type of arthritis caused by the buildup of uric acid in the blood, leading to sudden and severe pain, redness, swelling, and tenderness in joints, often in the big toe.",
				},
				{
					Name:        "Crohn's Disease",
					Description: "Crohn's disease is a type of inflammatory bowel disease (IBD) causing inflammation in the digestive tract, leading to abdominal pain, diarrhea, fatigue, and weight loss.",
				},
				{
					Name:        "Ulcerative Colitis",
					Description: "Ulcerative colitis is an inflammatory bowel disease characterized by inflammation and ulcers in the colon and rectum, leading to symptoms like abdominal pain, diarrhea, rectal bleeding, and fatigue.",
				},
				{
					Name:        "Endometriosis",
					Description: "Endometriosis is a painful disorder where tissue similar to the lining of the uterus grows outside the uterus, causing pelvic pain, heavy menstrual bleeding, and fertility problems.",
				},
				{
					Name:        "Polycystic Ovary Syndrome (PCOS)",
					Description: "PCOS is a hormonal disorder common among women of reproductive age, causing irregular periods, excess androgen levels, cysts on the ovaries, and difficulties in conceiving.",
				},
				{
					Name:        "Multiple Sclerosis (MS)",
					Description: "MS is a chronic disease affecting the central nervous system, causing a wide range of symptoms including fatigue, numbness or weakness in limbs, vision problems, and difficulties with coordination and balance.",
				},
				{
					Name:        "Parkinson's Disease",
					Description: "Parkinson's disease is a progressive nervous system disorder affecting movement, causing tremors, stiffness, slowness of movement, and difficulties with balance and coordination.",
				},
				{
					Name:        "Bipolar Disorder",
					Description: "Bipolar disorder, also known as manic-depressive illness, is a mental health condition characterized by extreme mood swings that include emotional highs (mania or hypomania) and lows (depression).",
				},
				{
					Name:        "Schizophrenia",
					Description: "Schizophrenia is a chronic mental disorder characterized by distorted thinking, hallucinations, delusions, and abnormal social behavior. It may cause disruptions in daily functioning.",
				},
				{
					Name:        "Epilepsy",
					Description: "Epilepsy is a neurological disorder characterized by recurrent seizures, which can vary from brief and nearly undetectable to long periods of vigorous shaking.",
				},
				{
					Name:        "Sleep Apnea",
					Description: "Sleep apnea is a sleep disorder where breathing repeatedly stops and starts during sleep, leading to loud snoring, abrupt awakenings, and daytime sleepiness.",
				},
				{
					Name:        "Hepatitis",
					Description: "Hepatitis refers to liver inflammation, often caused by viral infections (such as hepatitis A, B, or C) or due to other factors like alcohol, drugs, or autoimmune diseases, leading to symptoms like jaundice, abdominal pain, fatigue, and nausea.",
				},
				{
					Name:        "Irritable Bowel Syndrome (IBS)",
					Description: "IBS is a gastrointestinal disorder characterized by abdominal pain, bloating, changes in bowel habits (diarrhea or constipation), and often associated with stress or certain foods.",
				},
				{
					Name:        "Chronic Migraine",
					Description: "Chronic migraine is a subtype of migraine characterized by experiencing headaches on 15 or more days per month for at least three months, with at least eight days being migraines.",
				},
				{
					Name:        "Lupus",
					Description: "Lupus is a chronic autoimmune disease that can affect various parts of the body, causing inflammation, pain, and damage to joints, skin, kidneys, blood cells, heart, and lungs.",
				},
				{
					Name:        "Attention-Deficit/Hyperactivity Disorder (ADHD)",
					Description: "ADHD is a neurodevelopmental disorder characterized by difficulty paying attention, hyperactivity, impulsivity, and often interferes with daily functioning or development.",
				},
				{
					Name:        "Ovarian Cancer",
					Description: "Ovarian cancer begins in the ovaries and is often referred to as the 'silent killer' due to its subtle or absent symptoms in the early stages. Symptoms may include abdominal bloating, pelvic pain, and difficulty eating or feeling full quickly.",
				},
				{
					Name:        "Testicular Cancer",
					Description: "Testicular cancer affects the testicles and is most common in young or middle-aged men. Symptoms may include a lump or swelling in the testicle, pain, discomfort, or changes in size or shape of the testicles.",
				},
				{
					Name:        "Pancreatitis",
					Description: "Pancreatitis is inflammation of the pancreas, causing severe abdominal pain, nausea, vomiting, and in severe cases, complications like organ damage and infections.",
				},
				{
					Name:        "Cushing's Syndrome",
					Description: "Cushing's syndrome occurs due to prolonged exposure to high levels of cortisol, leading to symptoms such as weight gain, fatigue, muscle weakness, high blood pressure, and changes in skin appearance.",
				},
				{
					Name:        "Anemia",
					Description: "Anemia is a condition characterized by a lack of healthy red blood cells or hemoglobin in the blood, resulting in symptoms such as fatigue, weakness, pale skin, and shortness of breath.",
				},
			}

			db.Create(illnesses)
		}
	}

	return nil
}
