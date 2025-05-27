import NextAuth from "next-auth";
import DiscordProvider from "next-auth/providers/discord"

console.log(process.env.DISCORD_CLIENT_ID);

const handler = NextAuth({
  providers: [
    DiscordProvider({
      clientId: process.env.DISCORD_CLIENT_ID!,
      clientSecret: process.env.DISCORD_CLIENT_SECRET!,
    }),
  ],
  secret: process.env.NEXTAUTH_SECRET!,
});

export const GET = handler;
export const POST = handler;
