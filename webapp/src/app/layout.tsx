import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Dream Journal",
  description: "Record and visualize your dreams",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
