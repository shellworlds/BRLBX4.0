import { handleAuth, handleLogin } from "@auth0/nextjs-auth0";

export const GET = handleAuth({
  login: handleLogin({
    returnTo: "/portal",
    authorizationParams: {
      ...(process.env.AUTH0_AUDIENCE
        ? { audience: process.env.AUTH0_AUDIENCE }
        : {}),
      scope: "openid profile email offline_access",
    },
  }),
});
