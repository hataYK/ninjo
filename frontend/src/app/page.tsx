"use client";

import { useRouter } from "next/navigation";
import { useAtomValue, useSetAtom } from "jotai";
import { Box, Button, Heading, Text, VStack } from "@chakra-ui/react";
import { isAuthenticatedAtom, userAtom } from "@/stores/auth";
import { useLogout } from "@/api/generated/ninjo";
import { useEffect } from "react";

export default function Home() {
  const router = useRouter();
  const isAuthenticated = useAtomValue(isAuthenticatedAtom);
  const user = useAtomValue(userAtom);
  const setUser = useSetAtom(userAtom);
  const logoutMutation = useLogout();

  useEffect(() => {
    if (!isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated || !user) {
    return null;
  }

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        setUser(null);
        router.push("/login");
      },
    });
  };

  return (
    <Box minH="100vh" p={8}>
      <VStack gap={4} align="start">
        <Heading size="xl">Ninjo</Heading>
        <Text>ようこそ、{user.display_name} さん</Text>
        <Button onClick={handleLogout} variant="outline" loading={logoutMutation.isPending}>
          ログアウト
        </Button>
      </VStack>
    </Box>
  );
}
