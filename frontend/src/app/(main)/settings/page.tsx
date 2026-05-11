"use client";

import {
  Box,
  Card,
  Heading,
  HStack,
  Text,
  VStack,
} from "@chakra-ui/react";

const DAYS = ["日", "月", "火", "水", "木", "金", "土"];

export default function SettingsPage() {
  return (
    <VStack gap={6} align="stretch">
      <Heading size="lg">設定</Heading>

      {/* アバター設定 */}
      <Card.Root>
        <Card.Header>
          <Heading size="md">アバター</Heading>
        </Card.Header>
        <Card.Body>
          <HStack gap={4} flexWrap="wrap">
            {Array.from({ length: 6 }, (_, i) => (
              <Box
                key={i}
                w="80px"
                h="80px"
                bg="gray.100"
                borderRadius="lg"
                display="flex"
                alignItems="center"
                justifyContent="center"
                cursor="pointer"
                _hover={{ bg: "gray.200" }}
                borderWidth="2px"
                borderColor={i === 0 ? "blue.400" : "transparent"}
              >
                <Text fontSize="xs" color="gray.500">
                  preset_{String(i + 1).padStart(2, "0")}
                </Text>
              </Box>
            ))}
          </HStack>
        </Card.Body>
      </Card.Root>

      {/* 可処分時間設定 */}
      <Card.Root>
        <Card.Header>
          <Heading size="md">可処分時間</Heading>
        </Card.Header>
        <Card.Body>
          <VStack gap={3} align="stretch">
            <Text fontSize="sm" color="gray.600">
              曜日ごとの勉強可能時間を設定してください
            </Text>
            {DAYS.map((day, i) => (
              <HStack key={i} justify="space-between" py={1}>
                <Text fontSize="sm" w="40px">
                  {day}
                </Text>
                <Box flex="1" h="8px" bg="gray.100" borderRadius="full">
                  <Box h="8px" bg="blue.400" borderRadius="full" w="0%" />
                </Box>
                <Text fontSize="sm" color="gray.500" w="50px" textAlign="right">
                  0.0h
                </Text>
              </HStack>
            ))}
            <Text fontSize="xs" color="gray.400" textAlign="right">
              合計: 0.0h / 週
            </Text>
          </VStack>
        </Card.Body>
      </Card.Root>
    </VStack>
  );
}
