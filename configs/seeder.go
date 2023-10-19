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
	if !checkEnumIsExist(db, "") {
		db.Raw(
			`CREATE TYPE allergens_enum 
			AS ENUM (
				'SEDAN',
				'HATCHBACK',
				'MINIVAN'
			);`,
		)
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
	if db.Migrator().HasTable(&models.Allergy{}) {
		if err := db.First(&models.Allergy{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			allergies := []*models.Allergy{
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
			}

			db.Create(allergies)
		}
	}

	return nil
}
