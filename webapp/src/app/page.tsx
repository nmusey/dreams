import Image from "next/image";

export default function Home() {
  const { status } = useSession();

  useEffect(() => {
    if (status !== 'authenticated') {
      signIn('google');
    }
  }, [status]);

  return (
    <h1>Welcome!</h1>
  );
}
