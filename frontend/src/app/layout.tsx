import type { Metadata } from "next";
import { Playfair_Display } from "next/font/google";
import "./globals.css";
import "./layout.css";
import patternImg from "../../public/pattern.png"
import { BasicPage } from "@/components/basic-page/basic-page";

const playfairDisplay = Playfair_Display({
  subsets: ["latin"],
  variable: "--font-serif",
});

export const metadata: Metadata = {
  title: "Wishlist App",
  description: "blueetfm's personal wishlist app!",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const bgStyle = {
    backgroundColor: '#FFF9EE',
    backgroundImage: `url('${patternImg.src}')`,
    backgroundRepeat: 'repeat',
    backgroundSize: '1000px',
    minHeight: '100vh'
  }

  return (
    <html lang="en">
      <head>
        <script src="https://accounts.google.com/gsi/client" async></script>
      </head>
      <body style={bgStyle} className={`app-body ${playfairDisplay.variable}`}>
        <BasicPage>
          {children}
        </BasicPage>
      </body>
    </html>
  );
}
