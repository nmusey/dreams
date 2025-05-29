import Providers from "./providers";
import AuthWrapper from "@/app/components/AuthWrapper";
import DreamForm from "@/app/components/DreamForm";
import DreamList from "@/app/components/DreamList";

export default function Home() {
  return (
    <Providers>
      <AuthWrapper>
        <div className="flex flex-col items-center justify-center h-screen">
          <h1 className="text-4xl font-bold mb-8">Dreams App</h1>
          <DreamForm />
          <DreamList />
        </div>
      </AuthWrapper>
    </Providers>
  );
}
