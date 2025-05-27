import NextAuth from "next-auth";
import Providers from "next-auth/providers";

export default NextAuth({
  providers: [
    Providers.Google({
      clientId: process.env.GOOGLE_CLIENT_ID,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET,
    }),
  ],
  database: AppDataSource,
  jwt: {
    secret: process.env.JWT_SECRET,
  },
  callbacks: {
    async signIn(user, account, profile) {
      return true;
    },
    async redirect(url, baseUrl) {
      return url.startsWith(baseUrl) ? url : "/api/auth/unauthorized";
    },
  },
});
