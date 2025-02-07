
interface SidebarProps {
  title: string;
  children?: React.ReactNode;
}

export function Sidebar({ title, children }: SidebarProps) {
  return (
    <div className="h-full w-2/12">
      <div className="w-full flex flex-col gap-2">
        <h1 className="w-full text-2xl text-white font-bold">{title}</h1>
        <div className="border border-l-white border-r-0 border-t-0 border-b-0 border-opacity-15 pl-5">
          {children}
        </div>
      </div>
    </div>
  );
}
