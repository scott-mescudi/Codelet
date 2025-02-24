const LinesBackground = () => {
	return (
    <>
      <div className="w-full absolute h-full">
        <div className="w-full h-full -z-10 flex flex-row">
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
        </div>
      </div>
    </>
  );
}
export default function Home() {
  return (
    <>
      <div className="w-full flex flex-col items-center">
        <div className="h-[75vh] w-2/3 relative flex flex-col justify-center border border-t-0 border-r-0 border-white/5">
          <LinesBackground />
        </div>
        <div className="h-[75vh] w-2/3 relative flex flex-col justify-center border border-t-0 border-r-0 border-white/5">
          <LinesBackground />
        </div>
      </div>
    </>
  );
}
