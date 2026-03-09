export enum FieldTypes {
    TEXT = "text",
    RICH_TEXT = "rich_text",
    NUMBER = "number",
    DATE = "date",
    BOOLEAN = "boolean",
    RELATIONSHIP = "relationship",
    MEDIA = "media",
    TONES = "tones",
    LOCATION = "location",
    JSON = "json",
}

export const FieldTypeMetadata = {
    [FieldTypes.TEXT]: { label: "Text", icon: "Type", description: "Small or long text like title or description" },
    [FieldTypes.RICH_TEXT]: {
        label: "Rich Text",
        icon: "FileText",
        description: "A rich text editor with formatting options",
    },
    [FieldTypes.NUMBER]: { label: "Number", icon: "Hash", description: "Numbers (integer, float, decimal)" },
    [FieldTypes.DATE]: { label: "Date", icon: "Calendar", description: "A date picker with optional time" },
    [FieldTypes.BOOLEAN]: { label: "Boolean", icon: "CheckSquare", description: "Yes or No, 1 or 0, true or false" },
    [FieldTypes.RELATIONSHIP]: { label: "Relationship", icon: "Link", description: "Connect to another content type" },
    [FieldTypes.MEDIA]: { label: "Media", icon: "Image", description: "Images, videos, or documents" },
    [FieldTypes.TONES]: { label: "Tones", icon: "Palette", description: "Color picker or brand tones" },
    [FieldTypes.LOCATION]: { label: "Location", icon: "MapPin", description: "Latitude, longitude and address" },
    [FieldTypes.JSON]: { label: "JSON", icon: "Code", description: "Data in JSON format" },
};
