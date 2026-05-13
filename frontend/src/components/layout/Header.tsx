"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAtomValue, useSetAtom } from "jotai";
import { Box, Flex, Heading, HStack, Text, Button } from "@chakra-ui/react";
import { isAuthenticatedAtom, userAtom } from "@/stores/auth";
import { useLogout } from "@/api/generated/ninjo";

export function Header() {
  const pathname = usePathname();
  const router = useRouter();
  const isAuthenticated = useAtomValue(isAuthenticatedAtom);
  const user = useAtomValue(userAtom);
  const setUser = useSetAtom(userAtom);
  const logoutMutation = useLogout();

  if (!isAuthenticated) return null;

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        setUser(null);
        router.push("/login");
      },
    });
  };

  const navItems = [
    { href: "/", label: "ホーム" },
    { href: "/plans", label: "計画" },
    { href: "/settings", label: "設定" },
  ];

  return (
    <Box as="header" bg="white" borderBottomWidth="1px" borderColor="gray.200" px={6} py={3}>
      <Flex justify="space-between" align="center" maxW="1200px" mx="auto">
        <HStack gap={8}>
          <Link href="/">
            <Heading size="md" color="blue.600">Ninjo</Heading>
          </Link>
          <HStack as="nav" gap={4}>
            {navItems.map((item) => (
              <Link key={item.href} href={item.href}>
                <Text
                  fontWeight={pathname === item.href ? "bold" : "normal"}
                  color={pathname === item.href ? "blue.600" : "gray.600"}
                  _hover={{ color: "blue.500" }}
                  fontSize="sm"
                >
                  {item.label}
                </Text>
              </Link>
            ))}
          </HStack>
        </HStack>
        <HStack gap={4}>
          <Text fontSize="sm" color="gray.600">
            {user?.display_name}
          </Text>
          <Button
            size="sm"
            variant="ghost"
            onClick={handleLogout}
            loading={logoutMutation.isPending}
          >
            ログアウト
          </Button>
        </HStack>
      </Flex>
    </Box>
  );
}
