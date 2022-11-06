export interface ISmeGetRateRequest {
  //subproductId and UW_role will be provided by Munich Re as const
  subproductId: number;
  UW_role: string;
  company: Company;
  answers: {
    step1: SmeAnswerPolicy[];
    step2: SmeAnswerLocation[];
  };
}

export interface Company {
  vatnumber: { value: string };
  //turnover
  opre_eur: { value: number };
  employees: { value: number };
}

export interface SmeAnswerPolicy {
  slug: SlugPolicyCoverages;
  value: InsuranceAnswerValue | LegalDefenceAnswerValue;
}

export interface SmeAnswerLocation {
  //progressive id
  buildingId: string;
  value: BuildingAnswerValue;
}

export interface BuildingAnswerValue {
  ateco: Ateco;
  postcode: string;
  province: string;
  buildingType:
    | "masonry","reinforcedConcrete","antiSeismicLaminatedTimber","steel";
  // numberOfFloors and constructionYear are mandatory (together with buildingType) if earthquake cover is given
  numberOfFloors?: "ground_floor" | "first" | "second" | "greater_than_second";
  constructionYear?: "before1972" | "1972between2009" | "after2009";
  //alarm is mandatory if theft cover is given
  alarm?: "yes" | "no";
  typeOfInsurance: TypeOfInsuranceAnswerValue;
  answer: SmeAnswerLocationInsurance[];
}

export interface SmeAnswerLocationInsurance {
  slug: SlugLocationCoverages;
  value:
    | InsuranceAnswerValue
    | AssistanceAnswerValue
    | TypeOfInsuranceAnswerValue;
}

export interface Deductible {
  deductible?: number;
  // selfinsurance expecting a percentage e.g. 50% = 0.5
  // the case "20% min 5,000" is made with "deductbile=5,000 and selfInsurance=0.2"
  selfInsurance?: number;
}

export interface InsuranceAnswerValue extends Deductible {
  sumInsuredLimitOfIndemnity: number;
  typeOfSumInsured: "replacementValue" | "firstLoss";
}

export interface LegalDefenceAnswerValue {
  legalDefence: "basic" | "extended";
}

export interface AssistanceAnswerValue {
  assistance: "yes" | "no";
}

export interface TypeOfInsuranceAnswerValue {
  typeOfInsurance: "allRisks" | "namedPerils";
}

enum SlugPolicyCoverages {
  "third-party-liability" = "third-party-liability",
  "damage-to-goods-in-custody" = "damage-to-goods-in-custody",
  "defect-liability-workmanships" = "defect-liability-workmanships",
  "defect-liability-12-months" = "defect-liability-12-months",
  "defect-liability-dm-37-2008" = "defect-liability-dm-37-2008",
  "property-damage-due-to-theft" = "property-damage-due-to-theft",
  "damage-to-goods-course-of-works" = "damage-to-goods-course-of-works",
  "employers-liability" = "employers-liability",
  "product-liability" = "product-liability",
  "third-party-liability-construction-company" = "third-party-liability-construction-company",
  "legal-defence" = "legal-defence",
  "cyber" = "cyber",
}

enum SlugLocationCoverages {
  "building" = "building",
  "content" = "content",
  "lease-holders-interest" = "lease-holders-interest",
  "burst-pipe" = "burst-pipe",
  "power-surge" = "power-surge",
  "atmospheric-event" = "atmospheric-event",
  "sociopolitical-event" = "sociopolitical-event",
  "terrorism" = "terrorism",
  "earthquake" = "earthquake",
  "river-flood" = "river-flood",
  "water-damage" = "water-damage",
  "glass" = "glass",
  "machinery-breakdown" = "machinery-breakdown",
  "third-party-recourse" = "third-party-recourse",
  "theft" = "theft",
  "valuables-in-safe-strongrooms" = "valuables-in-safe-strongrooms",
  "valuables" = "valuables",
  "electronic-equipment" = "electronic-equipment",
  "increased-cost-of-working" = "increased-cost-of-working",
  "restoration-of-data" = "restoration-of-data",
  "software-under-license" = "software-under-license",
  "business-interruption" = "business-interruption",
  "property-owners-liability" = "property-owners-liability",
  "environmental-liability" = "environmental-liability",
  "assistance" = "assistance",
}

enum Ateco {
  "000001" = "000001",
  "000002" = "000002",
}
