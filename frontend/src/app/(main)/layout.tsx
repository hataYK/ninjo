"use client";

import { Box } from "@chakra-ui/react";
import { Header } from "@/components/layout/Header";
import { AuthGuard } from "@/components/layout/AuthGuard";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthGuard>
      <Box minH="100vh" bg="gray.50">
        <Header />
        <Box as="main" maxW="1200px" mx="auto" p={6}>
          {children}
        </Box>
      </Box>
    </AuthGuard>
  );
}
