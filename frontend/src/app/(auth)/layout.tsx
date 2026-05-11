"use client";

import { Box, Container, Heading, Text, VStack } from "@chakra-ui/react";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <Box minH="100vh" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
      <Container maxW="sm">
        <VStack gap={6}>
          <VStack gap={1}>
            <Heading size="2xl">Ninjo</Heading>
            <Text color="gray.600" fontSize="sm">
              勉強をやり切る
            </Text>
          </VStack>
          {children}
        </VStack>
      </Container>
    </Box>
  );
}
