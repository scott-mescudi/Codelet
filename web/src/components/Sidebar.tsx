
interface SidebarProps {
  title: string;
  key?: any
  children?: React.ReactNode;
}

export function Sidebar({ title, children }: SidebarProps) {
  return (
    <div className=" w-full select-none">
      <div className="w-full flex flex-col gap-2">
        <h1 className="w-full text-xl text-white font-bold">{title}</h1>
        <div className="border border-l-white border-r-0 border-t-0 border-b-0 border-opacity-15 flex flex-col ">
          {children}
        </div>
      </div>
    </div>
  );
}
