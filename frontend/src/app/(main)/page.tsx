"use client";

import { Box, Heading, Text, VStack, HStack, Card } from "@chakra-ui/react";
import { useAtomValue } from "jotai";
import { userAtom } from "@/stores/auth";

export default function HomePage() {
  const user = useAtomValue(userAtom);

  return (
    <VStack gap={6} align="stretch">
      <Heading size="lg">
        ようこそ、{user?.display_name} さん
      </Heading>

      <HStack gap={6} align="start" flexWrap="wrap">
        {/* アバターエリア */}
        <Card.Root flex="1" minW="280px">
          <Card.Header>
            <Heading size="md">アバター</Heading>
          </Card.Header>
          <Card.Body>
            <VStack gap={3}>
              <Box
                w="120px"
                h="120px"
                bg="gray.200"
                borderRadius="full"
                display="flex"
                alignItems="center"
                justifyContent="center"
              >
                <Text fontSize="sm" color="gray.500">
                  未設定
                </Text>
              </Box>
              <Text fontSize="sm" color="gray.500">
                設定画面でアバターを選択できます
              </Text>
            </VStack>
          </Card.Body>
        </Card.Root>

        {/* 今日のタスク */}
        <Card.Root flex="2" minW="320px">
          <Card.Header>
            <Heading size="md">今日のタスク</Heading>
          </Card.Header>
          <Card.Body>
            <Text color="gray.500" fontSize="sm">
              計画を作成すると、今日のタスクがここに表示されます
            </Text>
          </Card.Body>
        </Card.Root>
      </HStack>

      {/* 計画一覧 */}
      <Card.Root>
        <Card.Header>
          <HStack justify="space-between">
            <Heading size="md">学習計画</Heading>
          </HStack>
        </Card.Header>
        <Card.Body>
          <Text color="gray.500" fontSize="sm">
            まだ計画がありません。「計画」ページから新しい計画を作成しましょう。
          </Text>
        </Card.Body>
      </Card.Root>
    </VStack>
  );
}
