import { parse, format as formatDate, type Locale, formatDistance } from "date-fns";
import { es } from "date-fns/locale";

type LocaleData = {
    months: string[];
    monthsShort: string[];
    weekdays: string[];
    weekdaysShort: string[];
    weekdaysMin: string[];
};

const locales: Record<string, LocaleData> = {
    en: {
        months: [
            "January",
            "February",
            "March",
            "April",
            "May",
            "June",
            "July",
            "August",
            "September",
            "October",
            "November",
            "December",
        ],
        monthsShort: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
        weekdays: ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"],
        weekdaysShort: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
        weekdaysMin: ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"],
    },
    es: {
        months: [
            "Enero",
            "Febrero",
            "Marzo",
            "Abril",
            "Mayo",
            "Junio",
            "Julio",
            "Agosto",
            "Septiembre",
            "Octubre",
            "Noviembre",
            "Diciembre",
        ],
        monthsShort: ["Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"],
        weekdays: ["Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"],
        weekdaysShort: ["Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"],
        weekdaysMin: ["Do", "Lu", "Ma", "Mi", "Ju", "Vi", "Sa"],
    },
    pt: {
        months: [
            "Janeiro",
            "Fevereiro",
            "Março",
            "Abril",
            "Maio",
            "Junho",
            "Julho",
            "Agosto",
            "Setembro",
            "Outubro",
            "Novembro",
            "Dezembro",
        ],
        monthsShort: ["Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"],
        weekdays: ["Domingo", "Segunda-feira", "Terça-feira", "Quarta-feira", "Quinta-feira", "Sexta-feira", "Sábado"],
        weekdaysShort: ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"],
        weekdaysMin: ["Do", "Se", "Te", "Qa", "Qi", "Sx", "Sa"],
    },
};

export class DateTime extends Date implements Date {
    protected _date: Date;

    public constructor();
    public constructor(
        year?: number,
        month?: number,
        date?: number,
        hours?: number,
        minutes?: number,
        seconds?: number,
        ms?: number,
    );
    public constructor(arg: Date | string | number);
    public constructor(...args: any[]) {
        if (args.length === 0) {
            const now = Date.now();
            super(now);
            this._date = new Date(now);
        } else if (
            args.length === 1 &&
            (typeof args[0] === "string" || typeof args[0] === "number" || args[0] instanceof Date)
        ) {
            let timestamp: number;

            if (typeof args[0] === "string") {
                const input = args[0].trim();

                const isoLike = /^\d{4}-\d{2}-\d{2}(?:[ T]\d{2}:\d{2}(?::\d{2}(?:\.\d{3})?)?)?$/;
                if (isoLike.test(input)) {
                    timestamp = new Date(input + "Z").getTime();
                } else {
                    timestamp = new Date(input).getTime();
                }
            } else {
                timestamp = new Date(args[0]).getTime();
            }

            super(timestamp);
            this._date = new Date(timestamp);
        } else if (args.every((arg) => typeof arg === "number") && args.length >= 1 && args.length <= 7) {
            const [year, month = 0, day = 1, hours = 0, minutes = 0, seconds = 0, ms = 0] = args;

            const utc = Date.UTC(year, month, day, hours, minutes, seconds, ms);
            super(utc);
            this._date = new Date(utc);
        } else {
            throw new Error("Invalid DateTime initialization arguments");
        }

        if (!this._date?.getTime()) {
            throw new Error("Invalid DateTime initialization arguments");
        }
    }

    public getDate(): number {
        return super.getUTCDate();
    }

    public getDay(): number {
        return super.getUTCDay();
    }

    public getFullYear(): number {
        return super.getUTCFullYear();
    }

    public getHours(): number {
        return super.getUTCHours();
    }

    public getMilliseconds(): number {
        return super.getUTCMilliseconds();
    }

    public getMinutes(): number {
        return super.getUTCMinutes();
    }

    public getMonth(): number {
        return super.getUTCMonth();
    }

    public getSeconds(): number {
        return super.getUTCSeconds();
    }

    public getTime(): number {
        return super.getTime();
    }

    public getTimezoneOffset(): number {
        return 0;
    }

    public setDate(date: number): number {
        return super.setUTCDate(date);
    }

    public setFullYear(year: number, month?: number, date?: number): number {
        return super.setUTCFullYear(year, month ?? this._date.getUTCMonth(), date ?? this._date.getUTCDate());
    }

    public setHours(hours: number, min?: number, sec?: number, ms?: number): number {
        return super.setUTCHours(
            hours,
            min ?? this._date.getUTCMinutes(),
            sec ?? this._date.getUTCSeconds(),
            ms ?? this._date.getUTCMilliseconds(),
        );
    }

    public setMilliseconds(ms: number): number {
        return super.setUTCMilliseconds(ms);
    }

    public setMinutes(min: number, sec?: number, ms?: number): number {
        return super.setUTCMinutes(min, sec ?? this._date.getUTCSeconds(), ms ?? this._date.getUTCMilliseconds());
    }

    public setMonth(month: number, date?: number): number {
        return super.setUTCMonth(month, date ?? this._date.getUTCDate());
    }

    public setSeconds(sec: number, ms?: number): number {
        return super.setUTCSeconds(sec, ms ?? this._date.getUTCMilliseconds());
    }

    public setTime(time: number): number {
        return super.setTime(time);
    }

    public toDateString(): string {
        const days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
        const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

        const dayName = days[this.getDay()];
        const monthName = months[this.getMonth()];
        const day = this.getDate().toString().padStart(2, "0");
        const year = this.getFullYear();

        return `${dayName} ${monthName} ${day} ${year}`;
    }

    public toISOString(): string {
        return super.toISOString();
    }

    public toJSON(key?: any): string {
        return this.toISOString();
    }

    toLocaleDateString(): string;
    toLocaleDateString(locales?: string | string[], options?: Intl.DateTimeFormatOptions): string;
    toLocaleDateString(locales?: any, options?: any): string {
        return super.toLocaleDateString(locales, options);
    }

    toLocaleString(): string;
    toLocaleString(locales?: string | string[], options?: Intl.DateTimeFormatOptions): string;
    toLocaleString(locales?: any, options?: any): string {
        return super.toLocaleString(locales, options);
    }

    toLocaleTimeString(): string;
    toLocaleTimeString(locales?: string | string[], options?: Intl.DateTimeFormatOptions): string;
    toLocaleTimeString(locales?: any, options?: any): string {
        return super.toLocaleTimeString(locales, options);
    }

    public getUTCMonthName(): string {
        const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        return months[this.getMonth()];
    }

    public override toString(): string {
        return super.toString();
    }

    public override toTimeString(): string {
        const hh = this.getUTCHours().toString().padStart(2, "0");
        const mm = this.getUTCMinutes().toString().padStart(2, "0");
        const ss = this.getUTCSeconds().toString().padStart(2, "0");
        return `${hh}:${mm}:${ss} GMT`;
    }

    // getVarDate: () => VarDate;

    // readonly prototype: Date;

    [Symbol.toPrimitive](hint: "default"): string;
    [Symbol.toPrimitive](hint: "string"): string;
    [Symbol.toPrimitive](hint: "number"): number;
    [Symbol.toPrimitive](hint: string): string | number {
        return this._date[Symbol.toPrimitive](hint);
    }

    // Helpers methods not part of the Date interface
    public get isToday(): boolean {
        const today = new DateTime();
        return (
            this.getFullYear() === today.getFullYear() &&
            this.getMonth() === today.getMonth() &&
            this.getDate() === today.getDate()
        );
    }

    public format(pattern: string, locale: keyof typeof locales = "es"): string {
        const _locale = locales[locale];
        const d = this;
        const pad = (n: number, size = 2) => n.toString().padStart(size, "0");

        const tokens: Record<string, () => string> = {
            yyyy: () => d.getFullYear().toString(),
            YYYY: () => d.getFullYear().toString(),
            YY: () => d.getFullYear().toString().slice(-2),

            MMMM: () => _locale.months[d.getMonth()],
            MMM: () => _locale.monthsShort[d.getMonth()],
            MM: () => pad(d.getMonth() + 1),
            M: () => (d.getMonth() + 1).toString(),

            DD: () => pad(d.getDate()),
            D: () => d.getDate().toString(),
            dddd: () => _locale.weekdays[d.getDay()],
            ddd: () => _locale.weekdaysShort[d.getDay()],
            dd: () => _locale.weekdaysMin[d.getDay()],

            HH: () => pad(d.getHours()),
            H: () => d.getHours().toString(),
            hh: () => pad(((d.getHours() + 11) % 12) + 1),
            h: () => (((d.getHours() + 11) % 12) + 1).toString(),

            mm: () => pad(d.getMinutes()),
            m: () => d.getMinutes().toString(),
            ss: () => pad(d.getSeconds()),
            s: () => d.getSeconds().toString(),

            A: () => (d.getHours() >= 12 ? "PM" : "AM"),
            a: () => (d.getHours() >= 12 ? "pm" : "am"),
        };

        const tokenRegex = new RegExp(
            Object.keys(tokens)
                .sort((a, b) => b.length - a.length)
                .join("|"),
            "g",
        );

        return pattern.replace(tokenRegex, (match) => {
            const fn = tokens[match];
            return fn ? fn() : match;
        });
    }
}

window.DateTime = DateTime;

export enum OffsetMode {
    Up = "up",
    Down = "down",
}
export enum TimezoneType {
    Local = "local",
    UTC = "utc",
}

export class Carbon extends window.DateTime {
    public static readonly TINYTIME_FORMAT = "HH:mm";
    public static readonly TIME_FORMAT = "HH:mm:ss";
    public static readonly DATE_FORMAT = "yyyy-MM-dd";
    public static readonly DATETIME_FORMAT = "yyyy-MM-dd HH:mm:ss";
    public static readonly FORMAT_12H = "hh:mm aaaa";
    public static readonly FORMAT_ES = "dd-MM-yyyy HH:mm";
    //yyyy-MM-dd'T'HH:mm:ssX
    public static readonly ISO_FORMAT = "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'";
    public static readonly ISO_ALTER_FORMAT = "yyyy-MM-dd'T'HH:mm:ss'Z'";
    public static readonly DATEES_FORMAT = "dd-MM-yyyy";

    public static _options: { offset: number; timezone: string; locale: Locale } = {
        offset: -4,
        timezone: "America/Santiago",
        locale: es,
    };

    private _timezoneType: TimezoneType;

    private constructor(value: Date | string | number, timezoneType: TimezoneType) {
        super(value);
        this._timezoneType = timezoneType;
    }

    /**
     * Set global options for offset and timezone.
     */
    public static setGlobalOptions(options: { offset: number; timezone: string; locale?: Locale }) {
        Object.assign(Carbon._options, options);
    }

    /**
     * Create a Carbon instance from a date, string, or number.
     */
    public static from(
        date: Date | string | number,
        formatString: string = Carbon.DATETIME_FORMAT,
        timezoneType: TimezoneType = TimezoneType.UTC,
    ): Carbon {
        if (date instanceof Carbon) {
            return date;
        }

        if (date instanceof DateTime) {
            return new Carbon(date.getTime(), timezoneType);
        }

        if (typeof date === "string") {
            const pdate = parse(date, formatString, new Date());
            const _date = new Date(
                Date.UTC(
                    pdate.getFullYear(),
                    pdate.getMonth(),
                    pdate.getDate(),
                    pdate.getHours(),
                    pdate.getMinutes(),
                    pdate.getSeconds(),
                    pdate.getMilliseconds(),
                ),
            );

            return new Carbon(_date, timezoneType);
        }

        if (date instanceof Date) {
            return new Carbon(
                Date.UTC(
                    date.getFullYear(),
                    date.getMonth(),
                    date.getDate(),
                    date.getHours(),
                    date.getMinutes(),
                    date.getSeconds(),
                    date.getMilliseconds(),
                ),
                timezoneType,
            );
        }

        return new Carbon(date, timezoneType);
    }

    /**
     * Create a Carbon instance from an ISO 8601 string. UTC
     */
    public static parseISO(date: string): Carbon {
        return Carbon.from(new DateTime(date), Carbon.DATETIME_FORMAT, TimezoneType.UTC);
    }

    /**
     * Create a Carbon instance configured as UTC.
     */
    public static fromUTC(date: Date | string | number, formatString: string = Carbon.DATETIME_FORMAT): Carbon {
        return Carbon.from(date, formatString, TimezoneType.UTC);
    }

    public static fromLocal(date: Date | string | number, formatString: string = Carbon.DATETIME_FORMAT): Carbon {
        return Carbon.from(date, formatString, TimezoneType.Local);
    }

    public static utcNow() {
        const instance = Carbon.fromUTC(DateTime.now());
        instance._timezoneType = TimezoneType.UTC;
        return instance;
    }

    public static localNow() {
        const instance = this.utcNow();
        const absint = Math.abs(Carbon._options.offset);
        Carbon._options.offset > 0 ? instance.addHours(absint) : instance.subHours(absint);

        return new Carbon(instance.getTime(), TimezoneType.Local);
    }

    /**
     * Returns the input type of the date ('utc' or 'local').
     */
    public get timezoneType(): TimezoneType {
        return this._timezoneType;
    }

    /**
     * Converts the current date to local time.
     * Returns a new Carbon instance.
     */
    public toLocal(): Carbon {
        if (this.timezoneType === TimezoneType.Local) {
            return new Carbon(this.getTime(), TimezoneType.Local);
        }

        return this.applyOffset();
    }

    /**
     * Converts the current date to UTC.
     * Returns a new Carbon instance.
     */
    public toUTC(): Carbon {
        if (this.timezoneType === TimezoneType.UTC) {
            return new Carbon(this.getTime(), TimezoneType.UTC);
        }

        return this.applyOffset();
    }

    public clone(): Carbon {
        return new Carbon(this.getTime(), this.timezoneType);
    }

    public setTimezoneType(timezoneType: TimezoneType) {
        this._timezoneType = timezoneType;
        return this;
    }

    public applyOffset(): Carbon {
        const instance = this.clone();
        const offsetHours = Carbon._options.offset;

        if (instance.timezoneType === TimezoneType.UTC) {
            instance.addHours(offsetHours);
            instance._timezoneType = TimezoneType.Local;
        } else if (instance.timezoneType === TimezoneType.Local) {
            instance.addHours(-offsetHours);
            instance._timezoneType = TimezoneType.UTC;
        }

        return instance;
    }

    public copy(): Carbon {
        return this.clone();
    }

    /**
     * Formats the current date to a string.
     */
    public format(formatString: string = Carbon.FORMAT_ES): string {
        return formatDate(this, formatString, {
            locale: Carbon._options.locale,
        });
    }

    public isSame(other: DateTime, format?: string): boolean {
        if (format) {
            return this.format(format) === formatDate(other, format);
        }

        return this.getTime() === other.getTime();
    }

    public formatRelativeTo(other: DateTime): string {
        return formatDistance(this, other, { locale: Carbon._options.locale });
    }

    public addYears(years: number) {
        this.setFullYear(this.getFullYear() + years);
        return this;
    }

    public addMonths(months: number) {
        this.setMonth(this.getMonth() + months);
        return this;
    }

    public addWeeks(weeks: number) {
        this.setDate(this.getDate() + weeks * 7);
        return this;
    }

    public addDays(days: number) {
        this.setDate(this.getDate() + days);
        return this;
    }

    public addHours(hours: number) {
        this.setHours(this.getHours() + hours);
        return this;
    }

    public addMinutes(minutes: number) {
        this.setMinutes(this.getMinutes() + minutes);
        return this;
    }

    public addSeconds(seconds: number) {
        this.setSeconds(this.getSeconds() + seconds);
        return this;
    }

    public addTime(hours: number, minutes: number, seconds: number = 0) {
        this.setHours(this.getHours() + hours);
        this.setMinutes(this.getMinutes() + minutes);
        this.setSeconds(this.getSeconds() + seconds);
        return this;
    }

    public subYears(years: number) {
        this.setFullYear(this.getFullYear() - years);
        return this;
    }

    public subMonths(months: number) {
        this.setMonth(this.getMonth() - months);
        return this;
    }

    public subWeeks(weeks: number) {
        this.setDate(this.getDate() - weeks * 7);
        return this;
    }

    public subDays(days: number) {
        this.setDate(this.getDate() - days);
        return this;
    }

    public subHours(hours: number) {
        this.setHours(this.getHours() - hours);
        return this;
    }

    public subMinutes(minutes: number) {
        this.setMinutes(this.getMinutes() - minutes);
        return this;
    }

    public subSeconds(seconds: number) {
        this.setSeconds(this.getSeconds() - seconds);
        return this;
    }

    public subTime(hours: number, minutes: number, seconds: number = 0) {
        this.setHours(this.getHours() - hours);
        this.setMinutes(this.getMinutes() - minutes);
        this.setSeconds(this.getSeconds() - seconds);
        return this;
    }

    public isAfter(other: DateTime): boolean {
        return this.getTime() > other.getTime();
    }

    public isAfterOrEqual(other: DateTime): boolean {
        return this.getTime() >= other.getTime();
    }

    public isBefore(other: DateTime): boolean {
        return this.getTime() < other.getTime();
    }

    public isBeforeOrEqual(other: DateTime): boolean {
        return this.getTime() <= other.getTime();
    }

    public diffInDays(other: Date): number {
        return Math.floor(Math.abs(this.getTime() - other.getTime()) / (24 * 60 * 60 * 1000));
    }

    public diffInHours(other: Date): number {
        return Math.floor(Math.abs(this.getTime() - other.getTime()) / (60 * 60 * 1000));
    }

    public diffInMinutes(other: Date): number {
        return Math.floor(Math.abs(this.getTime() - other.getTime()) / (60 * 1000));
    }

    public diffInSeconds(other: Date): number {
        return Math.floor(Math.abs(this.getTime() - other.getTime()) / 1000);
    }

    public startOfDay(): Carbon {
        const startOfDay = new Carbon(this.getTime(), this.timezoneType);
        startOfDay.setHours(0, 0, 0, 0);
        return startOfDay;
    }

    public endOfDay(): Carbon {
        const endOfDay = new Carbon(this.getTime(), this.timezoneType);
        endOfDay.setHours(23, 59, 59, 999);
        return endOfDay;
    }

    public startOfWeek(): Carbon {
        const startOfWeek = new Carbon(this.getTime(), this.timezoneType);
        startOfWeek.setDate(startOfWeek.getDate() - startOfWeek.getDay());
        return startOfWeek;
    }

    public endOfWeek(): Carbon {
        const endOfWeek = new Carbon(this.getTime(), this.timezoneType);
        endOfWeek.setDate(endOfWeek.getDate() + 6 - endOfWeek.getDay());
        return endOfWeek;
    }

    public startOfMonth(): Carbon {
        const startOfMonth = new Carbon(this.getTime(), this.timezoneType);
        startOfMonth.setDate(1);
        return startOfMonth;
    }

    public endOfMonth(): Carbon {
        const endOfMonth = new Carbon(this.getTime(), this.timezoneType);
        endOfMonth.setDate(1);
        endOfMonth.setMonth(endOfMonth.getMonth() + 1);
        return endOfMonth;
    }

    public between(start: Date, end: Date): boolean {
        const time = this.getTime();
        return time >= start.getTime() && time <= end.getTime();
    }

    public isWithinMinutes(before: number, after: number): boolean {
        const now = Date.now();
        const timestamp = this.getTime();

        const minTime = now - before * 60_000;
        const maxTime = now + after * 60_000;

        return timestamp >= minTime && timestamp <= maxTime;
    }

    public toDatetime(): DateTime {
        return new DateTime(this.getTime());
    }

    public toDate(): Date {
        return new DateTime(this.getTime());
    }

    public isPast(): boolean {
        return this.getTime() < Date.now();
    }

    public isFuture(): boolean {
        return this.getTime() > Date.now();
    }
}

export const createWeekDays = (date: DateTime): Carbon[] => {
    const startOfWeek = Carbon.from(date).startOfWeek();
    const endOfWeek = Carbon.from(date).endOfWeek();

    const days: Carbon[] = [];
    while (startOfWeek.getTime() <= endOfWeek.getTime()) {
        days.push(startOfWeek.clone());
        startOfWeek.addDays(1);
    }

    return days;
};

export const eachDayOfWeek = (
    date: DateTime,
    func: (day: Carbon) => { ref: Carbon; [key: string]: any },
): { ref: Carbon; [key: string]: any }[] => {
    const startOfWeek = Carbon.from(date).startOfWeek();
    const days: { ref: Carbon; [key: string]: any }[] = [];
    for (let i = 0; i < 7; i++) {
        const day = startOfWeek.clone().addDays(i);
        if (func) {
            days.push(func(day));
        } else {
            days.push({ ref: day });
        }
    }
    return days;
};
