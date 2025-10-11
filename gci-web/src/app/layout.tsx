import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Providers from "./providers";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "GCI - Classes & Tasks",
  description: "CRUD demo for GCI Test",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
          >
              <div className="min-h-screen max-w-4xl mx-auto p-6">
                  <h1 className="text-2xl font-semibold mb-4">GCI - Classes & Tasks</h1>
                  <Providers>{children}</Providers>
              </div>
        
      </body>
    </html>
  );
}
