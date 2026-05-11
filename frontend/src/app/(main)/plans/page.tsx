"use client";

import Link from "next/link";
import { Box, Button, Card, Heading, Text, VStack } from "@chakra-ui/react";

export default function PlansPage() {
  return (
    <VStack gap={6} align="stretch">
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Heading size="lg">学習計画</Heading>
        <Link href="/plans/new">
          <Button colorPalette="blue">新しい計画を作成</Button>
        </Link>
      </Box>

      <Card.Root>
        <Card.Body>
          <Text color="gray.500" fontSize="sm" textAlign="center" py={8}>
            まだ計画がありません。新しい計画を作成しましょう。
          </Text>
        </Card.Body>
      </Card.Root>
    </VStack>
  );
}
