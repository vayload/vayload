import React, { useState, useEffect, useRef } from "react";
import {
    Home,
    LayoutDashboard,
    Tags,
    Users,
    Shield,
    Settings,
    Plug,
    DatabaseBackup,
    Activity,
    LogOut,
    Images,
    FolderPen,
    Library,
    Megaphone,
    Globe,
    FileText,
    HelpCircle,
    Bell,
    Search,
    ChevronDown,
    Plus,
    MoreHorizontal,
    Filter,
    X,
    Check,
    AlertCircle,
    Clock,
    Menu,
    Puzzle,
    Save,
    Trash2,
    CreditCard,
    Lock,
    Mail,
    Smartphone,
    Server,
    Code,
    ChevronRight,
    Eye,
    Download,
    Upload,
    User,
} from "lucide-react";

const PROJECTS = [
    { id: 101, name: "E-commerce Main", slug: "ecommerce-main", locale: "en" },
    { id: 102, name: "Corporate Blog", slug: "corp-blog", locale: "es" },
    { id: 103, name: "Mobile App API", slug: "mobile-app", locale: "en" },
];

const NOTIFICATIONS = [
    {
        id: 1,
        title: "Export completed",
        body: "Your content export is ready to download.",
        type: "success",
        time: "2m ago",
        read: false,
    },
    {
        id: 2,
        title: "New team member",
        body: 'Carlos requested access to "E-commerce Main".',
        type: "info",
        time: "1h ago",
        read: false,
    },
    {
        id: 3,
        title: "Storage warning",
        body: "You have used 85% of your asset storage.",
        type: "warning",
        time: "1d ago",
        read: true,
    },
];

const CONTENT_TYPES = [
    { id: 1, name: "Article", apiID: "article", count: 124 },
    { id: 2, name: "Product", apiID: "product", count: 856 },
    { id: 3, name: "Category", apiID: "category", count: 12 },
];

const ENTRIES = [
    {
        id: 1,
        title: "Summer Collection Launch",
        contentType: "Article",
        status: "published",
        author: "Ana Admin",
        updated: "2h ago",
    },
    {
        id: 2,
        title: "Top 10 Trends",
        contentType: "Article",
        status: "draft",
        author: "Carlos Editor",
        updated: "5h ago",
    },
    {
        id: 3,
        title: "Blue T-Shirt",
        contentType: "Product",
        status: "published",
        author: "Ana Admin",
        updated: "1d ago",
    },
];

// --- UI COMPONENTS ---

const Modal = ({ isOpen, onClose, title, children, footer }) => {
    if (!isOpen) return null;
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
            <div className="bg-white rounded-xl shadow-2xl w-full max-w-md overflow-hidden transform transition-all">
                <div className="flex justify-between items-center p-4 border-b border-gray-100">
                    <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-gray-600 p-1 rounded-md hover:bg-gray-100"
                    >
                        <X size={20} />
                    </button>
                </div>
                <div className="p-4">{children}</div>
                {footer && <div className="p-4 bg-gray-50 border-t border-gray-100">{footer}</div>}
            </div>
        </div>
    );
};

const ScrollableTabs = ({ tabs, activeTab, onTabChange }: any) => {
    const scrollRef = useRef(null);

    return (
        <div className="relative border-b border-gray-200 bg-white px-6">
            <div ref={scrollRef} className="flex overflow-x-auto gap-8 scrollbar-hide -mb-px">
                {tabs.map((tab: any) => {
                    const isActive = activeTab === tab.id;
                    return (
                        <button
                            key={tab.id}
                            onClick={() => onTabChange(tab.id)}
                            className={`whitespace-nowrap py-4 text-sm font-medium border-b-2 transition-all duration-200 ${
                                isActive
                                    ? "border-indigo-600 text-indigo-600"
                                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                            }`}
                        >
                            <span className="flex items-center gap-2">
                                {tab.icon && <tab.icon size={16} />}
                                {tab.label}
                                {tab.count !== undefined && (
                                    <span
                                        className={`text-xs px-2 py-0.5 rounded-full ${isActive ? "bg-indigo-50 text-indigo-600" : "bg-gray-100 text-gray-500"}`}
                                    >
                                        {tab.count}
                                    </span>
                                )}
                            </span>
                        </button>
                    );
                })}
            </div>
        </div>
    );
};

const SectionHeader = ({ title, subtitle, actions, breadcrumbs }: any) => (
    <div className="flex flex-col gap-4 mb-6 px-8 pt-8">
        {breadcrumbs && (
            <div className="flex items-center gap-2 text-xs text-gray-500">
                <Home size={12} className="text-gray-400" />
                {breadcrumbs.map((crumb, i) => (
                    <React.Fragment key={i}>
                        <ChevronRight size={12} className="text-gray-300" />
                        <span className={i === breadcrumbs.length - 1 ? "text-gray-800 font-medium" : ""}>{crumb}</span>
                    </React.Fragment>
                ))}
            </div>
        )}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <div>
                <h2 className="text-2xl font-bold text-gray-900 tracking-tight">{title}</h2>
                {subtitle && <p className="text-gray-500 text-sm mt-1">{subtitle}</p>}
            </div>
            {actions && <div className="flex gap-3">{actions}</div>}
        </div>
    </div>
);

const Badge = ({ children, color = "gray" }) => {
    const colors = {
        gray: "bg-gray-100 text-gray-700",
        green: "bg-green-100 text-green-700",
        yellow: "bg-amber-100 text-amber-700",
        blue: "bg-blue-100 text-blue-700",
        red: "bg-red-100 text-red-700",
        purple: "bg-purple-100 text-purple-700",
    };
    return (
        <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${colors[color] || colors.gray}`}>
            {children}
        </span>
    );
};

const NotificationPanel = ({ isOpen, onClose }) => {
    if (!isOpen) return null;
    return (
        <>
            <div className="fixed inset-0 z-30" onClick={onClose}></div>
            <div className="absolute right-4 top-16 z-40 w-96 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200">
                <div className="flex items-center justify-between p-4 border-b border-gray-50 bg-gray-50/50">
                    <h3 className="font-semibold text-gray-900">Notifications</h3>
                    <button className="text-xs text-indigo-600 font-medium hover:text-indigo-700">
                        Mark all as read
                    </button>
                </div>
                <div className="max-h-[400px] overflow-y-auto">
                    {NOTIFICATIONS.map((notif) => (
                        <div
                            key={notif.id}
                            className={`p-4 border-b border-gray-50 hover:bg-gray-50 transition-colors cursor-pointer flex gap-3 ${!notif.read ? "bg-indigo-50/30" : ""}`}
                        >
                            <div
                                className={`mt-1 shrink-0 w-2 h-2 rounded-full ${!notif.read ? "bg-indigo-500" : "bg-transparent"}`}
                            />
                            <div>
                                <p className="text-sm font-medium text-gray-900">{notif.title}</p>
                                <p className="text-sm text-gray-500 mt-0.5 line-clamp-2">{notif.body}</p>
                                <p className="text-xs text-gray-400 mt-2">{notif.time}</p>
                            </div>
                        </div>
                    ))}
                </div>
                <div className="p-3 bg-gray-50 border-t border-gray-100 text-center">
                    <button className="text-sm text-gray-600 font-medium hover:text-gray-900">View History</button>
                </div>
            </div>
        </>
    );
};

// --- VIEWS ---

const DashboardView = () => (
    <div className="pb-8">
        <SectionHeader
            title="Overview"
            subtitle="Welcome back, Ana. Here's what's happening today."
            breadcrumbs={["Dashboard", "Overview"]}
        />

        <div className="px-8 space-y-8">
            {/* Stats Row */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {[
                    {
                        label: "Total Revenue",
                        value: "$45,231",
                        change: "+20.1%",
                        trend: "up",
                        icon: Activity,
                        color: "text-indigo-600 bg-indigo-50",
                    },
                    {
                        label: "Active Users",
                        value: "+2,350",
                        change: "+18.1%",
                        trend: "up",
                        icon: Users,
                        color: "text-green-600 bg-green-50",
                    },
                    {
                        label: "Entries",
                        value: "12,234",
                        change: "+19%",
                        trend: "up",
                        icon: FileText,
                        color: "text-blue-600 bg-blue-50",
                    },
                    {
                        label: "Avg. Response",
                        value: "24ms",
                        change: "-4%",
                        trend: "down",
                        icon: Clock,
                        color: "text-purple-600 bg-purple-50",
                    },
                ].map((stat, i) => (
                    <div
                        key={i}
                        className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition-all"
                    >
                        <div className="flex items-center justify-between mb-4">
                            <span className="text-sm font-medium text-gray-500">{stat.label}</span>
                            <div className={`p-2 rounded-lg ${stat.color}`}>
                                <stat.icon size={18} />
                            </div>
                        </div>
                        <div className="flex items-baseline gap-2">
                            <span className="text-2xl font-bold text-gray-900">{stat.value}</span>
                            <span
                                className={`text-xs font-medium ${stat.trend === "up" ? "text-green-600" : "text-gray-500"}`}
                            >
                                {stat.change}
                            </span>
                        </div>
                    </div>
                ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Chart Area */}
                <div className="lg:col-span-2 bg-white rounded-xl border border-gray-200 shadow-sm p-6">
                    <div className="flex justify-between items-center mb-6">
                        <h3 className="text-lg font-semibold text-gray-900">Content Traffic</h3>
                        <select className="text-sm border-gray-300 rounded-md text-gray-500 bg-gray-50 px-2 py-1">
                            <option>Last 30 days</option>
                            <option>Last 7 days</option>
                        </select>
                    </div>
                    <div className="h-64 flex items-end gap-2 px-2">
                        {[40, 70, 45, 90, 60, 80, 50, 40, 70, 45, 90, 60, 75, 50, 85].map((h, i) => (
                            <div
                                key={i}
                                className="flex-1 bg-indigo-50 hover:bg-indigo-500 transition-colors rounded-t-sm relative group cursor-pointer"
                                style={{ height: `${h}%` }}
                            >
                                <div className="opacity-0 group-hover:opacity-100 absolute bottom-full left-1/2 -translate-x-1/2 mb-2 bg-gray-900 text-white text-xs py-1 px-2 rounded whitespace-nowrap z-10 transition-opacity">
                                    {h * 120} views
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Recent Activity */}
                <div className="bg-white rounded-xl border border-gray-200 shadow-sm p-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-6">Activity Feed</h3>
                    <div className="space-y-6 relative before:absolute before:left-2 before:top-2 before:bottom-2 before:w-0.5 before:bg-gray-100">
                        {[
                            {
                                user: "Ana Admin",
                                action: "published",
                                target: "Summer Sale",
                                time: "10m ago",
                            },
                            {
                                user: "Carlos Editor",
                                action: "updated",
                                target: "About Us",
                                time: "2h ago",
                            },
                            {
                                user: "System",
                                action: "backup",
                                target: "Daily Backup",
                                time: "5h ago",
                            },
                            {
                                user: "Ana Admin",
                                action: "invited",
                                target: "New User",
                                time: "1d ago",
                            },
                        ].map((item, i) => (
                            <div key={i} className="relative pl-8">
                                <div className="absolute left-0 top-1.5 w-4 h-4 rounded-full border-2 border-white bg-indigo-500 shadow-sm z-10"></div>
                                <p className="text-sm text-gray-900">
                                    <span className="font-medium">{item.user}</span> {item.action}{" "}
                                    <span className="font-medium">{item.target}</span>
                                </p>
                                <p className="text-xs text-gray-500 mt-0.5">{item.time}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    </div>
);

const EntriesView = () => {
    const [activeTab, setActiveTab] = useState("all");

    return (
        <div className="flex flex-col h-full bg-gray-50/50">
            <SectionHeader
                title="Entries"
                subtitle="Manage your content across all collections."
                breadcrumbs={["Workspace", "Entries"]}
                actions={
                    <button className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 shadow-sm flex items-center gap-2">
                        <Plus size={16} /> Create New
                    </button>
                }
            />

            <ScrollableTabs
                activeTab={activeTab}
                onTabChange={setActiveTab}
                tabs={[
                    { id: "all", label: "All Entries", count: 1250 },
                    {
                        id: "published",
                        label: "Published",
                        count: 980,
                        icon: Check,
                    },
                    { id: "draft", label: "Drafts", count: 45, icon: FileText },
                    {
                        id: "archived",
                        label: "Archived",
                        count: 225,
                        icon: DatabaseBackup,
                    },
                    {
                        id: "scheduled",
                        label: "Scheduled",
                        count: 5,
                        icon: Clock,
                    },
                ]}
            />

            <div className="flex-1 p-8 overflow-y-auto">
                <div className="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden">
                    {/* Filters Bar */}
                    <div className="p-4 border-b border-gray-200 flex flex-wrap gap-4 justify-between items-center bg-white">
                        <div className="flex items-center gap-2">
                            <div className="relative">
                                <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                                <input
                                    type="text"
                                    placeholder="Search entries..."
                                    className="pl-9 pr-4 py-2 border border-gray-200 rounded-lg text-sm w-64 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
                                />
                            </div>
                            <button className="px-3 py-2 border border-gray-200 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-50 flex items-center gap-2">
                                <Filter size={16} /> Filters
                            </button>
                        </div>
                        <div className="text-sm text-gray-500">Showing 1-10 of 1250</div>
                    </div>

                    <table className="min-w-full divide-y divide-gray-100">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="w-12 px-6 py-3">
                                    <input
                                        type="checkbox"
                                        className="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                                    />
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Title
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Content Type
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Status
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                    Last Updated
                                </th>
                                <th className="px-6 py-3 text-right"></th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-100">
                            {[
                                {
                                    title: "Summer Collection 2024",
                                    type: "Campaign",
                                    status: "published",
                                    date: "2 hours ago",
                                    author: "Ana",
                                },
                                {
                                    title: "Getting Started with Next.js",
                                    type: "Blog Post",
                                    status: "draft",
                                    date: "5 hours ago",
                                    author: "Carlos",
                                },
                                {
                                    title: "Premium Cotton T-Shirt",
                                    type: "Product",
                                    status: "published",
                                    date: "1 day ago",
                                    author: "Ana",
                                },
                                {
                                    title: "Q4 Financial Report",
                                    type: "Report",
                                    status: "archived",
                                    date: "1 week ago",
                                    author: "System",
                                },
                                {
                                    title: "Black Friday Landing Page",
                                    type: "Page",
                                    status: "scheduled",
                                    date: "In 2 days",
                                    author: "Ana",
                                },
                            ].map((entry, i) => (
                                <tr key={i} className="hover:bg-gray-50/80 transition-colors group cursor-pointer">
                                    <td className="px-6 py-4">
                                        <input
                                            type="checkbox"
                                            className="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                                        />
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="text-sm font-medium text-gray-900">{entry.title}</div>
                                        <div className="text-xs text-gray-500">by {entry.author}</div>
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">{entry.type}</td>
                                    <td className="px-6 py-4">
                                        <Badge
                                            color={
                                                entry.status === "published"
                                                    ? "green"
                                                    : entry.status === "draft"
                                                      ? "yellow"
                                                      : entry.status === "archived"
                                                        ? "gray"
                                                        : "blue"
                                            }
                                        >
                                            {entry.status}
                                        </Badge>
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">{entry.date}</td>
                                    <td className="px-6 py-4 text-right">
                                        <button className="p-1 text-gray-400 hover:text-indigo-600 rounded opacity-0 group-hover:opacity-100 transition-all">
                                            <MoreHorizontal size={18} />
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>

                    <div className="p-4 border-t border-gray-200 bg-gray-50 flex justify-center">
                        <button className="text-sm text-gray-600 font-medium hover:text-indigo-600">Load More</button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const ContentTypesView = ({ onOpenModal }) => (
    <div className="pb-8">
        <SectionHeader
            title="Content Types"
            subtitle="Define the schema and structure of your content."
            breadcrumbs={["Workspace", "Content Types"]}
            actions={
                <button
                    onClick={() => onOpenModal("create-type")}
                    className="cursor-pointer px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 shadow-sm flex items-center gap-2"
                >
                    <Plus size={16} /> Create Type
                </button>
            }
        />

        <div className="px-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
                {
                    name: "Blog Post",
                    api: "blog_post",
                    fields: 12,
                    entries: 145,
                    icon: FileText,
                },
                {
                    name: "Product",
                    api: "product",
                    fields: 24,
                    entries: 850,
                    icon: Tags,
                },
                {
                    name: "Author",
                    api: "author",
                    fields: 5,
                    entries: 12,
                    icon: Users,
                },
                {
                    name: "Category",
                    api: "category",
                    fields: 3,
                    entries: 8,
                    icon: Filter,
                },
                {
                    name: "Page",
                    api: "page",
                    fields: 8,
                    entries: 24,
                    icon: Globe,
                },
                {
                    name: "SEO Settings",
                    api: "seo",
                    fields: 15,
                    entries: 1,
                    icon: Search,
                    type: "Single",
                },
            ].map((type, i) => (
                <div
                    key={i}
                    className="group bg-white p-6 rounded-xl border border-gray-200 hover:border-indigo-400 hover:shadow-md transition-all cursor-pointer relative overflow-hidden"
                >
                    <div className="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button className="text-gray-400 hover:text-indigo-600">
                            <Settings size={18} />
                        </button>
                    </div>
                    <div className="flex items-start gap-4 mb-4">
                        <div className="p-3 bg-indigo-50 text-indigo-600 rounded-lg group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                            <type.icon size={24} />
                        </div>
                        <div>
                            <h3 className="font-bold text-gray-900">{type.name}</h3>
                            <code className="text-xs text-gray-500 bg-gray-100 px-1.5 py-0.5 rounded mt-1 inline-block">
                                {type.api}
                            </code>
                        </div>
                    </div>
                    <div className="flex items-center justify-between text-sm text-gray-500 mt-6 pt-4 border-t border-gray-50">
                        <span>{type.fields} Fields</span>
                        <span className="flex items-center gap-1">
                            {type.entries} {type.type === "Single" ? "Instance" : "Entries"}
                        </span>
                    </div>
                </div>
            ))}
        </div>
    </div>
);

const MediaLibraryView = () => {
    const [activeTab, setActiveTab] = useState("all");

    return (
        <div className="flex flex-col h-full bg-gray-50/50">
            <SectionHeader
                title="Media Library"
                subtitle="Manage images, videos, and documents."
                breadcrumbs={["Workspace", "Assets"]}
                actions={
                    <div className="flex gap-2">
                        <button className="px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-50 shadow-sm flex items-center gap-2">
                            <Plus size={16} /> New Folder
                        </button>
                        <button className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 shadow-sm flex items-center gap-2">
                            <Upload size={16} /> Upload Asset
                        </button>
                    </div>
                }
            />

            <ScrollableTabs
                activeTab={activeTab}
                onTabChange={setActiveTab}
                tabs={[
                    { id: "all", label: "All Assets" },
                    { id: "images", label: "Images", icon: Images },
                    { id: "videos", label: "Videos", icon: Activity },
                    { id: "docs", label: "Documents", icon: FileText },
                ]}
            />

            <div className="flex-1 p-8 overflow-y-auto">
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
                    {/* Folder Mockups */}
                    {["Marketing", "Product Shots", "Avatars"].map((folder, i) => (
                        <div
                            key={`f-${i}`}
                            className="group aspect-square bg-indigo-50/50 border border-indigo-100 rounded-xl flex flex-col items-center justify-center cursor-pointer hover:bg-indigo-100/50 hover:border-indigo-300 transition-all"
                        >
                            <FolderPen size={40} className="text-indigo-400 mb-2 fill-current opacity-80" />
                            <span className="text-sm font-medium text-indigo-900">{folder}</span>
                            <span className="text-xs text-indigo-500">12 items</span>
                        </div>
                    ))}

                    {/* Image Mockups */}
                    {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((item) => (
                        <div
                            key={item}
                            className="group relative aspect-square bg-white border border-gray-200 rounded-xl overflow-hidden hover:shadow-md transition-all cursor-pointer"
                        >
                            <div className="absolute inset-0 flex items-center justify-center bg-gray-100 text-gray-300">
                                <Images size={32} />
                            </div>
                            {/* Overlay */}
                            <div className="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors flex flex-col justify-end p-3 opacity-0 group-hover:opacity-100">
                                <div className="flex justify-between items-end text-white">
                                    <div>
                                        <p className="text-xs font-medium truncate w-24">IMG_00{item}.jpg</p>
                                        <p className="text-[10px] opacity-80">1.2 MB</p>
                                    </div>
                                    <button className="p-1 hover:bg-white/20 rounded">
                                        <MoreHorizontal size={16} />
                                    </button>
                                </div>
                            </div>
                            <div className="absolute top-2 left-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                <input type="checkbox" className="rounded text-indigo-600 focus:ring-indigo-500" />
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

const IntegrationsView = () => {
    const [activeTab, setActiveTab] = useState("payment");

    return (
        <div className="flex flex-col h-full bg-gray-50/50">
            <SectionHeader
                title="Integrations Marketplace"
                subtitle="Connect your project with third-party services."
                breadcrumbs={["System", "Integrations"]}
            />

            <ScrollableTabs
                activeTab={activeTab}
                onTabChange={setActiveTab}
                tabs={[
                    { id: "all", label: "All Apps" },
                    { id: "payment", label: "Payment", icon: CreditCard },
                    { id: "analytics", label: "Analytics", icon: Activity },
                    { id: "comm", label: "Communications", icon: Mail },
                    { id: "storage", label: "Storage", icon: DatabaseBackup },
                ]}
            />

            <div className="flex-1 p-8 overflow-y-auto">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {[
                        {
                            name: "Stripe",
                            cat: "Payment",
                            desc: "Accept payments globally with the world’s best processing platform.",
                            installed: true,
                        },
                        {
                            name: "PayPal",
                            cat: "Payment",
                            desc: "Simple and secure payment solutions for everyone.",
                            installed: false,
                        },
                        {
                            name: "SendGrid",
                            cat: "Email",
                            desc: "Cloud-based email delivery platform.",
                            installed: true,
                        },
                        {
                            name: "AWS S3",
                            cat: "Storage",
                            desc: "Scalable storage in the cloud.",
                            installed: true,
                        },
                        {
                            name: "Google Analytics",
                            cat: "Analytics",
                            desc: "Get insights into your users and traffic.",
                            installed: false,
                        },
                        {
                            name: "Algolia",
                            cat: "Search",
                            desc: "Fast, reliable search for your content.",
                            installed: false,
                        },
                    ].map((app, i) => (
                        <div
                            key={i}
                            className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm flex flex-col justify-between hover:border-indigo-400 transition-colors"
                        >
                            <div>
                                <div className="flex justify-between items-start mb-4">
                                    <div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center text-gray-500 font-bold text-lg">
                                        {app.name[0]}
                                    </div>
                                    {app.installed && <Badge color="green">Installed</Badge>}
                                </div>
                                <h4 className="font-bold text-gray-900">{app.name}</h4>
                                <p className="text-sm text-gray-500 mt-2">{app.desc}</p>
                            </div>
                            <button
                                className={`mt-6 w-full py-2 rounded-lg text-sm font-medium border transition-colors ${
                                    app.installed
                                        ? "border-gray-200 text-gray-700 hover:bg-gray-50"
                                        : "border-transparent bg-indigo-600 text-white hover:bg-indigo-700"
                                }`}
                            >
                                {app.installed ? "Configure" : "Install"}
                            </button>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

// --- UPDATED SETTINGS VIEW (HORIZONTAL TABS) ---

const SettingsView = () => {
    const [activeTab, setActiveTab] = useState("general");

    return (
        <div className="flex flex-col h-full bg-gray-50/50">
            <SectionHeader
                title="Project Settings"
                subtitle="Manage configuration, security, and billing."
                breadcrumbs={["System", "Settings"]}
                actions={
                    <button className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 shadow-sm flex items-center gap-2">
                        <Save size={16} /> Save Changes
                    </button>
                }
            />

            <ScrollableTabs
                activeTab={activeTab}
                onTabChange={setActiveTab}
                tabs={[
                    { id: "general", label: "General Details", icon: Settings },
                    { id: "security", label: "Security & Auth", icon: Lock },
                    {
                        id: "billing",
                        label: "Plan & Billing",
                        icon: CreditCard,
                    },
                    { id: "team", label: "Team Members", icon: Users },
                    { id: "api", label: "API Keys", icon: Code },
                    { id: "email", label: "Email Templates", icon: Mail },
                ]}
            />

            <div className="flex-1 p-8 overflow-y-auto">
                <div className="max-w-4xl mx-auto bg-white p-8 rounded-xl border border-gray-200 shadow-sm">
                    <h2 className="text-xl font-bold text-gray-900 mb-6 capitalize">{activeTab} Configuration</h2>

                    {/* Dynamic Form Content */}
                    <div className="space-y-6">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Project Name</label>
                                <input
                                    type="text"
                                    defaultValue="E-commerce Main"
                                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Project ID</label>
                                <input
                                    type="text"
                                    defaultValue="proj_829301"
                                    disabled
                                    className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-500"
                                />
                            </div>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Primary Domain</label>
                            <div className="flex">
                                <span className="inline-flex items-center px-3 rounded-l-lg border border-r-0 border-gray-300 bg-gray-50 text-gray-500 text-sm">
                                    https://
                                </span>
                                <input
                                    type="text"
                                    defaultValue="api.ecommerce.com"
                                    className="flex-1 px-4 py-2 border border-gray-300 rounded-r-lg focus:ring-2 focus:ring-indigo-500 outline-none"
                                />
                            </div>
                            <p className="mt-2 text-xs text-gray-500">
                                This domain will be used for your public API endpoints.
                            </p>
                        </div>

                        <div className="pt-6 border-t border-gray-100">
                            <h3 className="text-sm font-medium text-gray-900 mb-4">Danger Zone</h3>
                            <div className="flex items-center justify-between p-4 bg-red-50 border border-red-100 rounded-lg">
                                <div>
                                    <p className="text-sm font-medium text-red-800">Delete Project</p>
                                    <p className="text-xs text-red-600 mt-1">This action cannot be undone.</p>
                                </div>
                                <button className="px-4 py-2 bg-white border border-red-200 text-red-600 hover:bg-red-50 rounded-lg text-sm font-medium">
                                    Delete
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

const UsersView = () => (
    <div className="pb-8">
        <SectionHeader
            title="User Management"
            subtitle="Control access to your project."
            breadcrumbs={["System", "Users"]}
            actions={
                <button className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700">
                    Invite User
                </button>
            }
        />
        <div className="px-8">
            <div className="bg-white border border-gray-200 rounded-xl shadow-sm">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">User</th>
                            <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Role</th>
                            <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">
                                Status
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">
                                Last Login
                            </th>
                            <th className="px-6 py-3"></th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                        {[
                            {
                                name: "Ana Admin",
                                email: "ana@company.com",
                                role: "Owner",
                                status: "Active",
                                login: "2 mins ago",
                            },
                            {
                                name: "Carlos Dev",
                                email: "carlos@company.com",
                                role: "Developer",
                                status: "Active",
                                login: "1 day ago",
                            },
                            {
                                name: "Sarah Editor",
                                email: "sarah@company.com",
                                role: "Editor",
                                status: "Invited",
                                login: "-",
                            },
                        ].map((user, i) => (
                            <tr key={i} className="hover:bg-gray-50">
                                <td className="px-6 py-4">
                                    <div className="flex items-center">
                                        <div className="w-8 h-8 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center text-xs font-bold mr-3">
                                            {user.name[0]}
                                        </div>
                                        <div>
                                            <div className="text-sm font-medium text-gray-900">{user.name}</div>
                                            <div className="text-xs text-gray-500">{user.email}</div>
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-600">{user.role}</td>
                                <td className="px-6 py-4">
                                    <Badge color={user.status === "Active" ? "green" : "yellow"}>{user.status}</Badge>
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-500">{user.login}</td>
                                <td className="px-6 py-4 text-right">
                                    <button className="text-gray-400 hover:text-indigo-600">
                                        <MoreHorizontal size={18} />
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
);

// --- NEW ROLES VIEW ---

const RolesView = () => (
    <div className="pb-8">
        <SectionHeader
            title="Roles & Permissions"
            subtitle="Define what users can do in the system."
            breadcrumbs={["System", "Roles & ACL"]}
            actions={
                <button className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700">
                    Create Role
                </button>
            }
        />
        <div className="px-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
                {
                    name: "Administrator",
                    users: 3,
                    desc: "Full access to all system features.",
                    type: "System",
                },
                {
                    name: "Editor",
                    users: 12,
                    desc: "Can manage content but not system settings.",
                    type: "Custom",
                },
                {
                    name: "Author",
                    users: 5,
                    desc: "Can create and edit their own content.",
                    type: "Custom",
                },
                {
                    name: "Public API",
                    users: 0,
                    desc: "Read-only access for public endpoints.",
                    type: "System",
                },
            ].map((role, i) => (
                <div key={i} className="bg-white p-6 rounded-xl border border-gray-200 hover:shadow-md transition-all">
                    <div className="flex justify-between items-start mb-4">
                        <div className="p-2 bg-gray-100 rounded-lg text-gray-600">
                            <Shield size={20} />
                        </div>
                        <Badge color={role.type === "System" ? "purple" : "blue"}>{role.type}</Badge>
                    </div>
                    <h3 className="font-bold text-gray-900">{role.name}</h3>
                    <p className="text-sm text-gray-500 mt-2 mb-6 h-10">{role.desc}</p>
                    <div className="flex items-center justify-between pt-4 border-t border-gray-50">
                        <div className="flex -space-x-2">
                            {[...Array(Math.min(role.users, 3))].map((_, u) => (
                                <div key={u} className="w-6 h-6 rounded-full bg-gray-200 border-2 border-white"></div>
                            ))}
                            {role.users > 3 && (
                                <div className="w-6 h-6 rounded-full bg-gray-100 border-2 border-white flex items-center justify-center text-[8px] text-gray-500">
                                    +{role.users - 3}
                                </div>
                            )}
                        </div>
                        <button className="text-sm font-medium text-indigo-600 hover:text-indigo-800">
                            Edit Permissions
                        </button>
                    </div>
                </div>
            ))}
        </div>
    </div>
);

// --- APP LAYOUT ---

const Sidebar = ({ activeView, setActiveView, isMobileOpen, setIsMobileOpen }) => {
    const [currentProject, setCurrentProject] = useState({
        id: 101,
        name: "E-commerce Main",
        locale: "en-US",
    });
    const [isDropdownOpen, setDropdownOpen] = useState(false);
    const dropdownRef = useRef(null);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setDropdownOpen(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const menuGroups = [
        {
            id: "platform",
            label: "Platform",
            items: [
                { id: "dashboard", label: "Dashboard", icon: Home },
                {
                    id: "content-types",
                    label: "Content Types",
                    icon: LayoutDashboard,
                },
                { id: "entries", label: "Content Entries", icon: FolderPen },
                { id: "assets", label: "Media Library", icon: Images },
            ],
        },
        {
            id: "system",
            label: "System",
            items: [
                { id: "users", label: "User Management", icon: Users },
                { id: "roles", label: "Roles & ACL", icon: Shield },
                { id: "integrations", label: "Integrations", icon: Plug },
                { id: "audit", label: "Audit Logs", icon: FileText },
                { id: "settings", label: "Settings", icon: Settings },
            ],
        },
    ];

    const projects = [
        { id: 101, name: "E-commerce Main", locale: "en-US" },
        { id: 102, name: "Corporate Blog", locale: "es-ES" },
        { id: 103, name: "Mobile API", locale: "en-GB" },
    ];

    return (
        <>
            {isMobileOpen && (
                <div className="fixed inset-0 bg-black/50 z-30 md:hidden" onClick={() => setIsMobileOpen(false)} />
            )}
            <aside
                className={`fixed inset-y-0 left-0 z-40 bg-[#0f172a] text-gray-400 flex flex-col transition-all duration-300 ${isMobileOpen ? "translate-x-0 w-64" : "-translate-x-full md:translate-x-0 md:w-[72px] lg:w-64"}`}
            >
                {/* FIXED: Project Switcher with Working Popover */}
                <div className="relative h-16 border-b border-gray-800 shrink-0" ref={dropdownRef}>
                    <button
                        onClick={() => setDropdownOpen(!isDropdownOpen)}
                        className="w-full h-full flex items-center justify-center lg:justify-start px-0 lg:px-4 hover:bg-white/5 transition-colors"
                    >
                        <div className="flex items-center gap-3 w-full justify-center lg:justify-start">
                            <div className="w-8 h-8 rounded bg-indigo-500 flex items-center justify-center text-white font-bold shrink-0 shadow-lg shadow-indigo-500/30">
                                {currentProject.name.substring(0, 2).toUpperCase()}
                            </div>
                            <div className="hidden lg:block text-left min-w-0 flex-1">
                                <div className="text-sm font-medium text-gray-200 truncate">{currentProject.name}</div>
                                <div className="text-[10px] uppercase tracking-wider text-gray-500">
                                    {currentProject.locale}
                                </div>
                            </div>
                            <ChevronDown
                                size={14}
                                className={`hidden lg:block transition-transform ${isDropdownOpen ? "rotate-180" : ""}`}
                            />
                        </div>
                    </button>

                    {/* Dropdown Content */}
                    {isDropdownOpen && (
                        <div className="absolute top-14 left-2 right-2 bg-gray-800 border border-gray-700 rounded-xl shadow-2xl z-50 overflow-hidden animate-in fade-in slide-in-from-top-2">
                            <div className="py-2">
                                {projects.map((proj) => (
                                    <button
                                        key={proj.id}
                                        onClick={() => {
                                            setCurrentProject(proj);
                                            setDropdownOpen(false);
                                        }}
                                        className="w-full text-left px-4 py-2.5 text-sm hover:bg-gray-700 flex items-center justify-between group transition-colors"
                                    >
                                        <div>
                                            <div
                                                className={`font-medium ${proj.id === currentProject.id ? "text-white" : "text-gray-400 group-hover:text-gray-200"}`}
                                            >
                                                {proj.name}
                                            </div>
                                            <div className="text-[10px] text-gray-500">{proj.locale}</div>
                                        </div>
                                        {proj.id === currentProject.id && (
                                            <Check size={14} className="text-indigo-400" />
                                        )}
                                    </button>
                                ))}
                            </div>
                            <div className="border-t border-gray-700 p-2">
                                <button className="w-full flex items-center gap-2 px-3 py-2 text-xs font-medium text-indigo-400 hover:bg-indigo-500/10 rounded-lg transition-colors">
                                    <Plus size={14} /> Create New Project
                                </button>
                            </div>
                        </div>
                    )}
                </div>

                {/* Navigation */}
                <div className="flex-1 overflow-y-auto py-6 space-y-8">
                    {menuGroups.map((group) => (
                        <div key={group.id} className="px-3">
                            <div className="hidden lg:block px-3 mb-2 text-xs font-semibold uppercase tracking-wider text-gray-600">
                                {group.label}
                            </div>
                            <div className="space-y-1">
                                {group.items.map((item) => {
                                    const isActive = activeView === item.id;
                                    return (
                                        <button
                                            key={item.id}
                                            onClick={() => {
                                                setActiveView(item.id);
                                                setIsMobileOpen(false);
                                            }}
                                            className={`group relative flex items-center lg:justify-start justify-center w-full p-2.5 rounded-lg transition-all duration-200 ${isActive ? "bg-indigo-500/10 text-indigo-400" : "hover:bg-white/5 hover:text-gray-200"}`}
                                            title={item.label}
                                        >
                                            <item.icon size={20} strokeWidth={isActive ? 2.5 : 2} />
                                            <span className="hidden lg:block ml-3 text-sm font-medium">
                                                {item.label}
                                            </span>
                                            {isActive && (
                                                <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-indigo-500 rounded-r-full lg:hidden" />
                                            )}
                                        </button>
                                    );
                                })}
                            </div>
                        </div>
                    ))}
                </div>

                {/* FIXED: Removed User Profile Footer from Sidebar */}
                <div className="p-4 border-t border-gray-800 lg:hidden">
                    {/* Only show on mobile sidebar, on desktop it's in header */}
                    <button className="flex items-center gap-3 w-full p-2 rounded-lg hover:bg-white/5 transition-colors text-gray-400">
                        <LogOut size={20} />
                        <span className="font-medium text-sm">Logout</span>
                    </button>
                </div>
            </aside>
        </>
    );
};

const Header = ({ title, onMenuClick, hasUnread, onNotificationsClick }) => (
    <header className="sticky top-0 z-20 bg-white/80 backdrop-blur-md border-b border-gray-200 h-16 px-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
            <button onClick={onMenuClick} className="md:hidden p-2 text-gray-600 rounded-lg hover:bg-gray-100">
                <Menu size={20} />
            </button>
            <h1 className="text-lg font-semibold text-gray-800 hidden md:block">{title}</h1>
        </div>

        <div className="flex items-center gap-4">
            {/* Global Search */}
            <div className="hidden md:flex items-center bg-gray-100 rounded-lg px-3 py-1.5 focus-within:ring-2 focus-within:ring-indigo-500/20 transition-all">
                <Search size={16} className="text-gray-400" />
                <input
                    type="text"
                    placeholder="Jump to..."
                    className="bg-transparent border-none text-sm ml-2 w-48 focus:outline-none text-gray-700"
                />
                <div className="text-[10px] font-mono text-gray-400 border border-gray-300 rounded px-1.5 ml-2">⌘K</div>
            </div>

            <div className="h-6 w-px bg-gray-200 hidden md:block"></div>

            {/* Notifications */}
            <button
                onClick={onNotificationsClick}
                className={`relative p-2 rounded-lg transition-colors ${hasUnread ? "text-indigo-600 bg-indigo-50" : "text-gray-500 hover:bg-gray-100"}`}
            >
                <Bell size={20} />
                {hasUnread && (
                    <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full ring-2 ring-white"></span>
                )}
            </button>

            {/* FIXED: User Profile Moved to Header */}
            <button className="flex items-center gap-3 pl-2 pr-1 py-1 rounded-full hover:bg-gray-50 border border-transparent hover:border-gray-200 transition-all ml-2">
                <div className="text-right hidden sm:block leading-tight">
                    <div className="text-sm font-medium text-gray-900">Ana Admin</div>
                    <div className="text-xs text-gray-500">Super Admin</div>
                </div>
                <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-purple-600 to-blue-500 flex items-center justify-center text-white text-xs font-bold shadow-sm ring-2 ring-white">
                    AA
                </div>
                <ChevronDown size={14} className="text-gray-400 hidden sm:block mr-1" />
            </button>
        </div>
    </header>
);

const App = () => {
    const [activeView, setActiveView] = useState("dashboard");
    const [isMobileOpen, setIsMobileOpen] = useState(false);
    const [isNotifOpen, setIsNotifOpen] = useState(false);
    const [activeModal, setActiveModal] = useState(null);

    const getPageTitle = (viewId) => {
        const map = {
            dashboard: "Dashboard",
            "content-types": "Content Type Builder",
            entries: "Content Entries",
            assets: "Media Library",
            users: "User Management",
            roles: "Roles & Permissions",
            integrations: "Integrations",
            settings: "Project Settings",
            audit: "Audit Logs",
        };
        return map[viewId] || "Dashboard";
    };

    const renderContent = () => {
        switch (activeView) {
            case "dashboard":
                return <DashboardView />;
            case "content-types":
                return <ContentTypesView onOpenModal={setActiveModal} />;
            case "entries":
                return <EntriesView />;
            case "assets":
                return <MediaLibraryView />;
            case "users":
                return <UsersView />;
            // FIXED: Added Roles View
            case "roles":
                return <RolesView />;
            case "integrations":
                return <IntegrationsView />;
            case "settings":
                return <SettingsView />;
            case "audit":
                return (
                    <div className="pb-8">
                        <SectionHeader
                            title="Audit Logs"
                            subtitle="Track all system activity."
                            breadcrumbs={["System", "Audit Logs"]}
                        />
                        <div className="px-8">
                            <div className="bg-white p-12 text-center border border-gray-200 rounded-xl text-gray-500">
                                Audit Logs Table Placeholder
                            </div>
                        </div>
                    </div>
                );
            default:
                return <div className="p-8">Page Not Found</div>;
        }
    };

    return (
        <div className="flex h-screen bg-[#f8fafc] font-sans text-gray-900 overflow-hidden">
            <Sidebar
                activeView={activeView}
                setActiveView={setActiveView}
                isMobileOpen={isMobileOpen}
                setIsMobileOpen={setIsMobileOpen}
            />
            <div className="flex-1 flex flex-col min-w-0 md:pl-[72px] lg:pl-32 transition-all duration-300">
                <div className="flex-1 flex flex-col min-w-0 md:pl-[72px] lg:pl-32 transition-all duration-300">
                    <Header
                        title={getPageTitle(activeView)}
                        onMenuClick={() => setIsMobileOpen(true)}
                        onNotificationsClick={() => setIsNotifOpen(!isNotifOpen)}
                        hasUnread={true}
                    />

                    <div className="relative flex-1 overflow-hidden">
                        <NotificationPanel isOpen={isNotifOpen} onClose={() => setIsNotifOpen(false)} />

                        <main className="h-full overflow-y-auto">{renderContent()}</main>
                    </div>
                </div>
            </div>

            {/* GLOBAL MODALS */}
            <Modal
                isOpen={activeModal === "create-type"}
                onClose={() => setActiveModal(null)}
                title="Create Content Type"
                footer={
                    <div className="flex justify-end gap-2">
                        <button
                            onClick={() => setActiveModal(null)}
                            className="px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 rounded-lg"
                        >
                            Cancel
                        </button>
                        <button
                            onClick={() => setActiveModal(null)}
                            className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm"
                        >
                            Create Type
                        </button>
                    </div>
                }
            >
                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Display Name</label>
                        <input
                            type="text"
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
                            placeholder="e.g. Blog Post"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">API ID</label>
                        <input
                            type="text"
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-500 font-mono text-sm"
                            placeholder="blog_post"
                            disabled
                            value="blog_post"
                        />
                        <p className="text-xs text-gray-500 mt-1">Generated automatically from display name.</p>
                    </div>
                    <div className="flex items-center gap-2 mt-2">
                        <input
                            type="checkbox"
                            id="single"
                            className="rounded text-indigo-600 focus:ring-indigo-500 border-gray-300"
                        />
                        <label htmlFor="single" className="text-sm text-gray-700">
                            This is a Single Type (e.g. Homepage)
                        </label>
                    </div>
                </div>
            </Modal>
        </div>
    );
};

export default App;
