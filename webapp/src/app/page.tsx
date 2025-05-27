import Providers from "./providers";
import AuthWrapper from "./components/AuthWrapper";

export default function Home() {
  return (
    <Providers>
      <AuthWrapper>
        <h1>Welcome!</h1>
      </AuthWrapper>
    </Providers>
  );
}
