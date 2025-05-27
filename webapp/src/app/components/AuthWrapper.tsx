"use client";

import { signIn, useSession } from "next-auth/react";
import { ReactNode, useEffect } from "react";

interface AuthWrapperProps {
    children: ReactNode;
}

export default function AuthWrapper({ children }: AuthWrapperProps) {
    const { status } = useSession() || "";
    useEffect(() => {
        if (status !== 'authenticated') {
            signIn('discord');
        }
    }, [status]);

    if (status !== 'authenticated') {
        return null;
    }

    return <>{children}</>;
} 